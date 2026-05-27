# AI/RAG Interview Cheatsheet

## Core RAG Concepts

| Term | Meaning | Simple Example |
| :--- | :--- | :--- |
| **RAG** (Retrieval-Augmented Generation) | Retrieve relevant documents first, then generate answer using them | User asks "How to refund?" → System searches knowledge base → Finds refund policy → LLM answers using that policy |
| **Embedding** | Numerical representation of text (vector of ~1536 numbers) | "cat" → [0.12, -0.45, 0.89, ...] (1536 numbers) |
| **Vector Database** | Database that stores embeddings and searches by similarity | Pinecone, Weaviate, OpenSearch – find "dog" near "puppy" in vector space |
| **Similarity Search** | Finding vectors closest to query vector using cosine similarity | Query "refund policy" returns most similar chunks from knowledge base |
| **Hybrid Search** | Combines keyword search (BM25) + semantic search (vectors) | "REF-12345" uses BM25; "how do I cancel" uses semantic; both combined |
| **BM25** | Keyword matching algorithm with term frequency weighting | Search "error code E104" → finds exact match in docs, ignores synonyms |
| **Cross-Encoder Reranker** | Second-stage model that scores query-document pairs together | After hybrid search returns 20 chunks, reranker scores each (query, chunk) pair to pick top 5 |
| **Chunking** | Splitting long documents into smaller pieces for retrieval | 10-page policy doc → split into 512-token chunks with 20% overlap |
| **Context Window** | Max tokens LLM can process at once | GPT-4o has 128k tokens; you need retrieved chunks + query + instructions to fit |

---

## Agentic AI Concepts

| Term | Meaning | Simple Example |
| :--- | :--- | :--- |
| **Agentic System** | AI that can take actions, use tools, make decisions | Customer support agent that decides: "I need to check order status" → calls order API → returns answer |
| **Agent Orchestrator** | Central agent that routes tasks to specialized sub-agents | Main agent: "User wants refund" → calls Refund Agent to check policy, calls Order Agent to verify purchase |
| **Tool Calling** | LLM decides to call an external function/API | User: "What's my order status?" → LLM decides to call `get_order_status(order_id)` |
| **Multi-Agent Architecture** | Multiple AI agents collaborating on tasks | Orchestrator + Retrieval Agent + Tool Agent + Validation Agent working together |
| **Memory/State Handling** | Agent remembers previous turns in conversation | User: "My name is Ahmed" → Agent stores in memory → Later: "What's my name?" → Agent remembers "Ahmed" |
| **Reflection Pattern** | Agent reviews its own output to improve | Agent generates answer → Second agent reviews for accuracy → Suggests improvements |
| **Structured Output** | LLM returns JSON/Pydantic instead of free text | Output: `{"intent": "refund", "confidence": 0.92, "order_id": "12345"}` |

---

## LLM & Generation Concepts

| Term | Meaning | Simple Example |
| :--- | :--- | :--- |
| **Temperature** | Controls randomness of LLM output (0 = deterministic, 1 = creative) | Temperature 0.2 for factual QA; Temperature 0.8 for brainstorming |
| **Token** | Smallest unit of text LLM processes (~4 chars in English) | "Hello world" → ["Hello", " world"] → 2 tokens |
| **Streaming** | Sending tokens to user as they're generated | User sees "The" then " answer" then " is..." instead of waiting 2 seconds |
| **Hallucination** | LLM generates false information confidently | Q: "When was Ozmo founded?" A: "2015" (actual date is 2018) |
| **Prompt Engineering** | Crafting instructions to get desired output | "You are a support agent. Only answer using provided context. If unsure, say 'I don't know'." |
| **Few-Shot Prompting** | Providing examples in the prompt | Q: "Cancel order" → A: "Refund processed" [example] → Then actual query |
| **System Prompt** | Permanent instructions that don't change per query | "You are Ozmo support. Be helpful, concise, and cite sources." |
| **User Prompt** | The actual user query | "How do I reset my password?" |
| **Logit Probabilities** | Raw probability scores for each token before selection | Model 60% sure next token is "refund", 30% "cancel", 10% "return" → confidence 0.6 |
| **Context Engineering** | Structuring what information is sent to LLM | System instructions + retrieved chunks + conversation history + user query |

---

## Guardrails & Validation

| Term | Meaning | Simple Example |
| :--- | :--- | :--- |
| **Guardrails** | Safety checks before/after LLM generation | Block query containing "hack my account"; block response containing "credit card number" |
| **Hallucination Detection** | Checking if LLM response contradicts retrieved chunks | LLM says "refund takes 3 days" but policy says "5-7 days" → flag hallucination |
| **Confidence Score** | Numerical measure of how certain the system is | 0.92 = very confident; 0.45 = not confident → escalate to human |
| **PII Filter** | Remove personal identifiable information | "My email is ahmed@gmail.com" → redacted to "[EMAIL]" |
| **Toxicity Filter** | Block offensive content | Response contains profanity → reject and return fallback message |
| **Citation Validation** | Verify LLM claims are actually in retrieved chunks | LLM says "policy says X" → check if X exists in retrieved chunk → if not, hallucination |
| **Human-in-the-Loop (HITL)** | Human reviews/approves AI output before action | LLM drafts response → human agent approves or edits → then sent to user |
| **Fallback Strategy** | What happens when LLM fails or low confidence | Confidence < 0.7 → route to human; LLM API down → use cached responses |

---

## Evaluation & Observability

| Term | Meaning | Simple Example |
| :--- | :--- | :--- |
| **Faithfulness** | Does answer stay true to retrieved context? | Context says "refund takes 5 days" → Answer "refund takes 5 days" = faithful (1.0) |
| **Answer Relevance** | Does answer actually address the question? | Q: "How to refund?" → A: "Our products are great" = low relevance (0.2) |
| **Context Recall** | Did retrieval get ALL needed information? | User asks about refund policy and exceptions → Context only has policy, not exceptions = partial recall |
| **RAGAS** | Framework for evaluating RAG systems | Computes faithfulness, answer relevance, context recall automatically |
| **Drift Detection** | Monitoring when model or data changes over time | Confidence scores dropped from 0.85 → 0.65 over 2 weeks → knowledge base may be stale |
| **Offline Evaluation** | Testing on historical data before deployment | Run 500 old queries through new model, compare to ground truth answers |
| **A/B Testing** | Compare two versions with live traffic | 50% users see model A, 50% see model B → measure escalation rates |

---

## Human Escalation

| Term | Meaning | Simple Example |
| :--- | :--- | :--- |
| **Escalation** | Transferring conversation from AI to human | Confidence 0.45 → "Let me connect you to a support agent" |
| **Draft + Approve** | AI drafts response, human approves/edits | LLM writes answer → Agent dashboard shows it → Agent clicks "Send" or edits |
| **Session Transfer** | Human takes over entire conversation | Agent clicks "Take Over" → all future messages go to human, not LLM |
| **Agent Queue** | FIFO queue of conversations waiting for human | SQS FIFO queue with 3 escalations waiting, each with conversation context |
| **Co-pilot Mode** | Human types, LLM suggests completions in real-time | Agent typing "Your refund will..." → LLM suggests "be processed in 5-7 days" → Agent accepts |

---

## Infrastructure (AWS)

| Term | Meaning | Simple Example |
| :--- | :--- | :--- |
| **API Gateway** | Entry point for all client requests | POST /chat → routes to orchestrator service; rate limits 100 requests/second |
| **Cognito** | User authentication and tenant isolation | User logs in → gets JWT with `tenant_id: "acme-corp"` → authorizer validates |
| **ECS Fargate** | Serverless container orchestration | Run orchestrator service: 2 vCPU, 8GB RAM, auto-scales based on queue length |
| **Microservices** | Small, independent services for each concern | Orchestrator Service + Retrieval Service + LLM Gateway + Validation Service |
| **SQS** | Message queue for async processing | User feedback → SQS queue → Lambda processes without blocking |
| **SNS** | Pub/sub for broadcasting events | "Customer resolved" event → SNS → fans out to analytics, training, audit queues |
| **OpenTelemetry (ADOT)** | Distributed tracing standard | Trace ID travels from API Gateway → Orchestrator → LLM → see end-to-end latency |
| **X-Ray** | AWS tracing service (deprecated, use OTel) | Still common in interviews, but say "I'd use OpenTelemetry" |
| **ServiceLens** | CloudWatch dashboard combining traces, logs, metrics | One view showing service map + trace timelines + error logs |

---

## Key Frameworks & Tools (From Job Description)

| Term | Meaning | Simple Example |
| :--- | :--- | :--- |
| **LangChain** | Framework for building LLM applications | `chain = prompt | llm | output_parser` → run with user query |
| **CrewAI** | Framework for multi-agent collaboration | Researcher Agent + Writer Agent + Reviewer Agent working on one task |
| **PydanticAI** | Type-safe LLM responses with Pydantic | `class Response(BaseModel): answer: str; confidence: float` → LLM returns validated object |
| **Claude Code** | AI coding assistant for development | Prompt: "Write a function to validate emails" → Claude generates code in your IDE |
| **AI Coding Agent** | AI that writes code as part of development process | You write comment "// retry logic with backoff" → AI generates implementation |

---

## Production Concepts (From Job Description)

| Term | Meaning | Simple Example |
| :--- | :--- | :--- |
| **Non-deterministic Behavior** | LLM can give different answers for same input | Same query "What's 2+2?" → sometimes "4", sometimes "four" |
| **Deterministic Services** | Code that always produces same output for same input | `add(2,2)` always returns `4` |
| **Progressive Delivery** | Gradual rollout of changes | 1% users → 10% → 50% → 100% while monitoring metrics |
| **Blast Radius** | Impact scope when something fails | New model affects only 1% of traffic → small blast radius |
| **Multi-Tenant SaaS** | Single system serves multiple customers securely | Tenant A and Tenant B share infrastructure but see only their own data |
| **Data Isolation** | Preventing cross-tenant data leaks | Every query includes `WHERE tenant_id = current_tenant_id` |
| **On-Call Support** | Engineers respond to production incidents | PagerDuty alerts when escalation rate > 20% → engineer investigates |
| **Incident Response** | Process for fixing production issues | Detect → Mitigate → Investigate → Fix → Post-mortem |

---

## Metrics & Tradeoffs (Mention These)

| Term | Meaning | Tradeoff |
| :--- | :--- | :--- |
| **Latency vs. Accuracy** | Faster vs. more correct answers | Faster = smaller model (worse accuracy); Accurate = slower (cross-encoder reranking) |
| **Cost vs. Quality** | Cheaper vs. better outputs | Cheap = GPT-3.5; Expensive = GPT-4o |
| **Precision vs. Recall** | Correctness vs. completeness of retrieval | High precision = fewer false positives; High recall = fewer missed docs |
| **Cache Hit Rate** | % of queries answered from cache | Higher = faster/cheaper; Lower = more fresh data |
| **Escalation Rate** | % of conversations sent to humans | Lower = AI handling more; Higher = safer but expensive |

---

## Sample Sentences to Say in Interview

> *"For hybrid search, I'd use α=0.7 favoring semantic for most queries, but classify product codes as BM25-only to ensure exact matches."*

> *"I'd validate context before LLM generation—check that retrieved chunks answer the question—then stream the response. Post-hoc validation runs async to catch hallucinations for training."*

> *"The orchestrator agent maintains conversation memory in DynamoDB and decides whether to call the retrieval agent, tool agent, or escalate to human based on confidence scores."*

> *"For multi-tenant isolation, every query includes tenant_id in the JWT claim, and I enforce it in vector DB queries and DynamoDB partition keys."*

> *"I'd use offline evaluation with RAGAS metrics (faithfulness, answer relevance) to block deployments that degrade quality by more than 5%."*

---

## Last-Minute Check

| Concept | Can you explain it in 1 sentence? |
| :--- | :--- |
| RAG | ✓ |
| Embedding | ✓ |
| Hybrid search | ✓ |
| Agent orchestrator | ✓ |
| Hallucination detection | ✓ |
| Human-in-the-loop | ✓ |
| Confidence threshold | ✓ |
| Multi-tenant isolation | ✓ |
| Progressive delivery | ✓ |

**Good luck Thursday! You've got this.**