class LLMResponder:
    """
    Send retrieved chunks + query to LLM and get structured response.
    """
    
    def __init__(self):
        self.llm = ChatOpenAI(
            model="gpt-4o",  # or "gpt-3.5-turbo" for cheaper
            temperature=0.2,  # Low temp = more factual, less creative
            streaming=True  # Enable streaming for better UX
        )
    
    def create_prompt_with_context(
        self, 
        query: str, 
        chunks: List[str], 
        conversation_history: List[Dict] = None
    ) -> str:
        """
        Build the prompt with system instructions, context, and conversation history.
        """
        # Combine chunks into context
        context = "\n\n".join([
            f"[Source {i+1}]: {chunk}" 
            for i, chunk in enumerate(chunks)
        ])
        
        # Build conversation history section
        history_text = ""
        if conversation_history:
            history_text = "Previous conversation:\n"
            for msg in conversation_history[-4:]:  # Last 4 messages only
                role = "User" if msg["role"] == "user" else "Assistant"
                history_text += f"{role}: {msg['content']}\n"
        
        # Complete prompt
        prompt = f"""
You are Ozmo Customer Support AI. Follow these rules STRICTLY:

1. ONLY use information from the provided context to answer.
2. If the context doesn't contain the answer, say: "I don't have that information. Let me connect you to a human agent."
3. ALWAYS cite your sources by adding [Source X] at the end of sentences.
4. Be concise, helpful, and professional.
5. Do NOT make up information or hallucinate.

{history_text}

Context from knowledge base:
{context}

User question: {query}

Your answer (with citations):
"""
        return prompt
    
    def get_response_with_confidence(
        self, 
        query: str, 
        chunks: List[str], 
        conversation_history: List[Dict] = None
    ) -> Dict:
        """
        Get response from LLM with confidence scoring.
        
        Returns:
            Dict with answer, confidence, citations, and whether to escalate
        """
        prompt = self.create_prompt_with_context(query, chunks, conversation_history)
        
        # Invoke LLM
        from langchain_core.messages import HumanMessage
        response = self.llm.invoke([HumanMessage(content=prompt)])
        
        answer = response.content
        
        # Simple confidence heuristic (in production, use logit probabilities)
        confidence = 0.8  # Default
        if "I don't have that information" in answer:
            confidence = 0.2
        elif "let me connect you" in answer.lower():
            confidence = 0.3
        elif len(chunks) == 0:
            confidence = 0.1
        
        # Extract citations (anything with [Source X])
        import re
        citations = re.findall(r'\[Source \d+\]', answer)
        
        return {
            "answer": answer,
            "confidence": confidence,
            "citations": citations,
            "should_escalate": confidence < 0.6,
            "chunks_used": chunks
        }
    
    async def stream_response(
        self, 
        query: str, 
        chunks: List[str]
    ):
        """
        Stream response token by token for better UX.
        
        This is async so you can send tokens to WebSocket as they arrive.
        """
        prompt = self.create_prompt_with_context(query, chunks)
        
        from langchain_core.messages import HumanMessage
        
        # Streaming version
        async for chunk in self.llm.astream([HumanMessage(content=prompt)]):
            yield chunk.content  # Send each token to the client

# Example usage:
# responder = LLMResponder()
# result = responder.get_response_with_confidence(
#     query="How long does a refund take?",
#     chunks=["Refunds take 5-7 business days to process [Source: Policy]"],
#     conversation_history=[{"role": "user", "content": "I want my money back"}]
# )
# print(f"Answer: {result['answer']}")
# print(f"Confidence: {result['confidence']}")
# if result['should_escalate']:
#     print("Escalating to human")