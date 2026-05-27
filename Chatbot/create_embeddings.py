def create_embeddings(texts: List[str]) -> List[List[float]]:
    """
    Convert text to vector embeddings using OpenAI's embedding model.
    
    Embeddings are numerical representations of text that capture semantic meaning.
    Similar text = similar vectors (close in vector space).
    
    Args:
        texts: List of text chunks to convert to embeddings
        
    Returns:
        List of embedding vectors (each is list of 1536 floats)
    """
    # Initialize OpenAI embeddings
    embeddings_model = OpenAIEmbeddings(
        model="text-embedding-3-small",  # Good balance of quality/cost
        dimensions=1536  # Default dimension for this model
    )
    
    # Generate embeddings in batch (more efficient than one-by-one)
    # embed_documents is for indexing (multiple texts)
    embeddings = embeddings_model.embed_documents(texts)
    
    # For a single query, use embed_query()
    # query_embedding = embeddings_model.embed_query("How to reset password?")
    
    return embeddings

# Example usage:
# docs = ["Refund policy: 5-7 business days", "Password reset: click forgot password"]
# vectors = create_embeddings(docs)
# print(f"Created {len(vectors)} embeddings, each with {len(vectors[0])} dimensions")