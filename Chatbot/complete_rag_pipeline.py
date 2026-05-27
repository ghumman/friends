def complete_rag_pipeline(
    user_query: str,
    knowledge_base_docs: List[str]
):
    """
    Complete RAG pipeline from query to safe response.
    
    This is what you'll describe in your interview.
    """
    
    # Step 1: Create embeddings (indexing - done once)
    # In production, this happens offline when documents are added
    retriever = HybridRetriever(knowledge_base_docs)
    
    # Step 2: Retrieve relevant chunks
    retrieved = retriever.classify_query_and_search(user_query)
    chunks = [r["text"] for r in retrieved]
    
    # Step 3: Validate before LLM
    validator = DataValidator()
    should_escalate, reason = validator.should_escalate_to_human(
        user_query, chunks, 0.75
    )
    
    if should_escalate:
        return {
            "answer": "I need to connect you to a human agent.",
            "escalate": True,
            "reason": reason
        }
    
    # Step 4: Generate response with guardrails
    pipeline = SafeLLMPipeline()
    result = pipeline.process(user_query, chunks)
    
    return result

# Example usage:
# docs = [
#     "Refund policy: Refunds are processed within 5-7 business days.",
#     "Password reset: Use the 'Forgot Password' link on the login page.",
#     "Shipping: Standard shipping takes 3-5 business days."
# ]
# 
# response = complete_rag_pipeline("How long for a refund?", docs)
# print(response["answer"])