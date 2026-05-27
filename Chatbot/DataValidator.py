class DataValidator:
    """
    Validate retrieved chunks BEFORE sending to LLM.
    
    This prevents wasting LLM calls on bad data.
    """
    
    def __init__(self):
        self.llm = ChatOpenAI(model="gpt-4o", temperature=0)
    
    def check_chunk_relevance(
        self, 
        query: str, 
        chunks: List[str]
    ) -> Tuple[bool, List[str], float]:
        """
        Check if retrieved chunks actually answer the user's question.
        
        Returns:
            (is_relevant, relevant_chunks, average_confidence)
        """
        if not chunks:
            return False, [], 0.0
        
        # Use LLM to evaluate relevance (fast with small prompt)
        prompt = f"""
        You are evaluating if retrieved document chunks can answer a user's question.
        
        User question: {query}
        
        Retrieved chunks:
        {chr(10).join([f"{i+1}. {chunk[:200]}..." for i, chunk in enumerate(chunks)])}
        
        Answer with JSON:
        {{
            "can_answer": true/false,
            "relevant_chunk_indices": [0, 1, 2],  # Which chunks are actually useful
            "confidence": 0.95  # How confident are you?
        }}
        
        Only return JSON, no other text.
        """
        
        from langchain_core.messages import HumanMessage
        response = self.llm.invoke([HumanMessage(content=prompt)])
        
        # Parse response (simplified - in production use Pydantic)
        import json
        try:
            result = json.loads(response.content)
            relevant_indices = result.get("relevant_chunk_indices", [])
            relevant_chunks = [chunks[i] for i in relevant_indices if i < len(chunks)]
            return (
                result.get("can_answer", False),
                relevant_chunks,
                result.get("confidence", 0.0)
            )
        except:
            # Fallback: return first chunk if parsing fails
            return bool(chunks), [chunks[0]], 0.5
    
    def check_query_safety(self, query: str) -> Tuple[bool, str]:
        """
        Check if query contains harmful content (PII, toxicity, injection).
        """
        # Simple PII detection (in production use regex + NLP)
        pii_patterns = [
            r'\b\d{3}[-.]?\d{3}[-.]?\d{4}\b',  # Phone number
            r'\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b',  # Email
            r'\b\d{16}\b'  # Credit card (simplified)
        ]
        
        import re
        for pattern in pii_patterns:
            if re.search(pattern, query):
                return False, f"Query contains PII: {re.search(pattern, query).group()}"
        
        # Injection detection (looking for attempts to override system prompts)
        injection_keywords = [
            "ignore previous instructions",
            "you are now", "system prompt", "forget your",
            "you are a different AI"
        ]
        
        query_lower = query.lower()
        for keyword in injection_keywords:
            if keyword in query_lower:
                return False, f"Query contains potential injection: '{keyword}'"
        
        return True, "Query is safe"
    
    def should_escalate_to_human(
        self, 
        query: str, 
        chunks: List[str], 
        relevance_confidence: float
    ) -> Tuple[bool, str]:
        """
        Decide if this query should go to a human instead of LLM.
        
        Escalate when:
        1. No relevant chunks found
        2. Very low confidence in retrieval
        3. Query contains harmful content
        """
        # Check safety first
        is_safe, reason = self.check_query_safety(query)
        if not is_safe:
            return True, f"Security: {reason}"
        
        # Check if we have any chunks
        if not chunks:
            return True, "No relevant documents found in knowledge base"
        
        # Check chunk relevance
        is_relevant, relevant_chunks, confidence = self.check_chunk_relevance(query, chunks)
        
        if not is_relevant:
            return True, f"Retrieved chunks don't answer the question (confidence: {confidence:.2f})"
        
        if confidence < 0.6:
            return True, f"Low confidence in retrieval quality: {confidence:.2f}"
        
        return False, "Ready for LLM generation"

# Example usage:
# validator = DataValidator()
# query = "Give me my refund now"
# chunks = ["Refund policy: takes 5-7 business days", "Contact support@ozmo.com"]
# escalate, reason = validator.should_escalate_to_human(query, chunks, 0.72)
# if escalate:
#     print(f"Escalating: {reason}")
# else:
#     print("Proceeding to LLM")