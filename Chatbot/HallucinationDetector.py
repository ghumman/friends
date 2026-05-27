class HallucinationDetector:
    """
    Detect if LLM is making up information not in retrieved context.
    
    This runs AFTER the LLM generates the response (asynchronous).
    """
    
    def __init__(self):
        self.llm = ChatOpenAI(model="gpt-4o", temperature=0)
    
    def check_hallucination(
        self, 
        query: str, 
        answer: str, 
        retrieved_chunks: List[str]
    ) -> Tuple[bool, float, List[str]]:
        """
        Check if the answer contains hallucinated information.
        
        Returns:
            (is_hallucinated, confidence_score, hallucinated_claims)
        """
        # Join chunks into single context
        context = "\n".join(retrieved_chunks)
        
        prompt = f"""
You are a hallucination detection system. Compare the AI's answer against the source context.

Source context (knowledge base):
{context}

AI's answer to user question "{query}":
{answer}

Task: Identify ANY statement in the AI's answer that is NOT supported by the source context.

Output JSON:
{{
    "has_hallucination": true/false,
    "hallucinated_claims": ["claim 1", "claim 2"],
    "confidence": 0.95,
    "explanation": "Brief explanation"
}}

Only return JSON, no other text.
"""
        from langchain_core.messages import HumanMessage
        response = self.llm.invoke([HumanMessage(content=prompt)])
        
        import json
        try:
            result = json.loads(response.content)
            return (
                result.get("has_hallucination", False),
                result.get("confidence", 0.5),
                result.get("hallucinated_claims", [])
            )
        except:
            # If parsing fails, assume no hallucination (fail open)
            return False, 0.5, []
    
    def check_faithfulness(self, answer: str, chunks: List[str]) -> float:
        """
        Calculate faithfulness score (0-1). Higher = more faithful to source.
        
        Faithfulness = does answer stay true to retrieved context?
        """
        if not chunks:
            return 0.0
        
        context = "\n".join(chunks)
        
        prompt = f"""
Score how faithful this answer is to the source context (0.0 to 1.0).

1.0 = Every statement in answer is directly supported by context
0.5 = Partial support, some claims not in context  
0.0 = Answer completely contradicts or makes up information

Source context: {context}

Answer: {answer}

Return ONLY a number between 0.0 and 1.0, no other text.
"""
        from langchain_core.messages import HumanMessage
        response = self.llm.invoke([HumanMessage(content=prompt)])
        
        try:
            score = float(response.content.strip())
            return min(max(score, 0.0), 1.0)  # Clamp to 0-1
        except:
            return 0.5
    
    def auto_correct_hallucination(
        self, 
        query: str, 
        hallucinated_answer: str, 
        chunks: List[str]
    ) -> str:
        """
        Attempt to correct a hallucinated answer by forcing the LLM to use only context.
        
        Use this as a fallback before escalating to human.
        """
        context = "\n".join(chunks)
        
        correction_prompt = f"""
The previous answer contained incorrect information NOT found in the source context.

Source context:
{context}

User question: {query}

Previous (incorrect) answer: {hallucinated_answer}

Please provide a CORRECT answer using ONLY information from the source context.
If the context doesn't contain the answer, say: "I cannot answer that based on available information."

Correct answer:
"""
        from langchain_core.messages import HumanMessage
        response = self.llm.invoke([HumanMessage(content=correction_prompt)])
        
        return response.content

# Example usage:
# detector = HallucinationDetector()
# chunks = ["Refunds take 5-7 business days"]
# answer = "Refunds take 2-3 business days"  # Hallucinated!
# 
# is_hallucinated, conf, claims = detector.check_hallucination(
#     query="How long for refund?", 
#     answer=answer, 
#     retrieved_chunks=chunks
# )
# print(f"Hallucination detected: {is_hallucinated}")
# print(f"Claims: {claims}")
# 
# if is_hallucinated:
#     corrected = detector.auto_correct_hallucination(query, answer, chunks)
#     print(f"Corrected answer: {corrected}")