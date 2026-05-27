class HybridRetriever:
    """
    Combines keyword search (BM25) and semantic search (embeddings).
    
    Why hybrid?
    - BM25: Good for exact matches (product codes, error messages)
    - Semantic: Good for synonyms and conceptual queries ("cancel" vs "terminate")
    - Together: Best of both worlds
    """
    
    def __init__(self, documents: List[str]):
        """
        Initialize with documents to search over.
        
        Args:
            documents: List of text chunks from your knowledge base
        """
        self.documents = documents
        
        # Setup BM25 (keyword search)
        # Tokenize documents into words (simple whitespace split)
        tokenized_docs = [doc.lower().split() for doc in documents]
        self.bm25 = BM25Okapi(tokenized_docs)
        
        # Setup semantic search (embeddings)
        self.embeddings_model = OpenAIEmbeddings(model="text-embedding-3-small")
        # Pre-compute embeddings for all documents (do this once at index time)
        self.doc_embeddings = self.embeddings_model.embed_documents(documents)
    
    def semantic_search(self, query: str, k: int = 10) -> List[Tuple[str, float]]:
        """
        Find documents semantically similar to query using cosine similarity.
        
        Cosine similarity = 1.0 when vectors point same direction (semantically identical)
        Cosine similarity = 0.0 when vectors are perpendicular (unrelated)
        """
        # Get embedding for the query
        query_embedding = self.embeddings_model.embed_query(query)
        
        # Calculate cosine similarity between query and each document
        similarities = []
        for doc_emb in self.doc_embeddings:
            # Cosine similarity = dot product / (norm1 * norm2)
            dot_product = np.dot(query_embedding, doc_emb)
            norm1 = np.linalg.norm(query_embedding)
            norm2 = np.linalg.norm(doc_emb)
            similarity = dot_product / (norm1 * norm2) if norm1 * norm2 > 0 else 0
            similarities.append(similarity)
        
        # Get top k indices
        top_indices = np.argsort(similarities)[-k:][::-1]
        
        return [(self.documents[i], similarities[i]) for i in top_indices]
    
    def bm25_search(self, query: str, k: int = 10) -> List[Tuple[str, float]]:
        """
        Find documents using keyword matching (BM25 algorithm).
        
        BM25 gives higher scores to documents that contain query terms frequently,
        but penalizes common words (like "the", "a") and very long documents.
        """
        tokenized_query = query.lower().split()
        scores = self.bm25.get_scores(tokenized_query)
        
        # Get top k indices
        top_indices = np.argsort(scores)[-k:][::-1]
        
        return [(self.documents[i], scores[i]) for i in top_indices]
    
    def hybrid_search(
        self, 
        query: str, 
        alpha: float = 0.7, 
        k: int = 10
    ) -> List[Dict]:
        """
        Combine BM25 and semantic search with weighted scoring.
        
        Args:
            query: User's question
            alpha: Weight for semantic search (0 = BM25 only, 1 = semantic only)
            k: Number of results to return
            
        Returns:
            List of dicts with text, score, and source information
        """
        # Get results from both methods
        semantic_results = self.semantic_search(query, k=20)  # Get extra for reranking
        bm25_results = self.bm25_search(query, k=20)
        
        # Normalize scores to 0-1 range for fair combination
        semantic_scores = [score for _, score in semantic_results]
        bm25_scores = [score for _, score in bm25_results]
        
        # Avoid division by zero
        max_semantic = max(semantic_scores) if semantic_scores else 1
        max_bm25 = max(bm25_scores) if bm25_scores else 1
        
        # Create a dictionary to combine scores by document text
        combined_scores = {}
        
        for text, score in semantic_results:
            normalized = score / max_semantic if max_semantic > 0 else 0
            combined_scores[text] = alpha * normalized
        
        for text, score in bm25_results:
            normalized = score / max_bm25 if max_bm25 > 0 else 0
            if text in combined_scores:
                # Document appears in both result sets, add BM25 score with weight (1-alpha)
                combined_scores[text] += (1 - alpha) * normalized
            else:
                combined_scores[text] = (1 - alpha) * normalized
        
        # Sort by combined score and return top k
        sorted_results = sorted(
            combined_scores.items(), 
            key=lambda x: x[1], 
            reverse=True
        )[:k]
        
        return [
            {"text": text, "combined_score": score, "source": "hybrid"}
            for text, score in sorted_results
        ]
    
    def classify_query_and_search(self, query: str) -> List[Dict]:
        """
        Smart routing: use different strategies based on query type.
        
        This is more efficient than always using hybrid search.
        """
        query_lower = query.lower()
        
        # Rule 1: Product codes or error codes (e.g., "REF-12345", "E104")
        if any(char.isdigit() for char in query) and '-' in query:
            print("Detected product/error code → using BM25 only")
            results = self.bm25_search(query, k=5)
            return [{"text": text, "score": score, "source": "bm25"} 
                    for text, score in results]
        
        # Rule 2: Short, exact-match queries (likely FAQ titles)
        if len(query.split()) <= 3:
            print("Detected short query → weighting BM25 higher")
            return self.hybrid_search(query, alpha=0.3)  # Favor BM25
        
        # Rule 3: Long, conversational questions
        if len(query.split()) > 10:
            print("Detected long query → weighting semantic higher")
            return self.hybrid_search(query, alpha=0.8)  # Favor semantic
        
        # Default: balanced hybrid
        print("Using default hybrid search")
        return self.hybrid_search(query, alpha=0.7)

# Example usage:
# docs = [
#     "Refund policy: Refunds take 5-7 business days to process",
#     "Product code REF-12345 is for the wireless mouse",
#     "Password reset: Click 'Forgot password' on login screen"
# ]
# retriever = HybridRetriever(docs)
# results = retriever.classify_query_and_search("How to get my money back?")
# for r in results:
#     print(f"Score: {r['score']:.3f} | Text: {r['text'][:50]}...")