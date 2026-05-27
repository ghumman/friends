class Guardrails:
    """
    Safety filters that run BEFORE and AFTER LLM processing.
    
    Input Guardrails: Block harmful queries before they reach LLM
    Output Guardrails: Block harmful responses before they reach user
    """
    
    def __init__(self):
        # Blocked topics (can be expanded with regex or ML)
        self.blocked_topics = [
            "hack", "exploit", "bypass security", "steal", "fraud",
            "illegal", "drug", "weapon", "terrorism", "discriminate",
            "racist", "sexist", "harass", "abuse"
        ]
        
        # Sensitive data patterns
        self.pii_patterns = {
            "ssn": r'\b\d{3}-\d{2}-\d{4}\b',
            "credit_card": r'\b\d{4}[-\s]?\d{4}[-\s]?\d{4}[-\s]?\d{4}\b',
            "email": r'\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b',
            "phone": r'\b\d{3}[-.]?\d{3}[-.]?\d{4}\b',
            "address": r'\b\d{1,5}\s[A-Za-z0-9\s,]+(Street|St|Avenue|Ave|Road|Rd|Boulevard|Blvd)\b',
        }
        
        # Redaction mapping
        self.redaction_map = {
            "ssn": "[SSN_REDACTED]",
            "credit_card": "[CREDIT_CARD_REDACTED]",
            "email": "[EMAIL_REDACTED]",
            "phone": "[PHONE_REDACTED]",
            "address": "[ADDRESS_REDACTED]"
        }
    
    def input_guardrail(self, query: str) -> Tuple[bool, str]:
        """
        Check if input query is safe to process.
        
        Returns:
            (is_allowed, rejection_reason_or_modified_query)
        """
        query_lower = query.lower()
        
        # Check for blocked topics
        for topic in self.blocked_topics:
            if topic in query_lower:
                return False, f"Query blocked: Contains prohibited content '{topic}'"
        
        # Check for prompt injection
        injection_patterns = [
            "ignore previous", "forget your", "you are now",
            "system prompt", "new instruction", "override",
            "pretend you are", "act as if", "roleplay"
        ]
        
        for pattern in injection_patterns:
            if pattern in query_lower:
                return False, f"Query blocked: Potential prompt injection detected '{pattern}'"
        
        # Check for excessive length (prevent DoS)
        if len(query) > 4000:
            return False, "Query too long (max 4000 characters)"
        
        return True, query
    
    def output_guardrail(
        self, 
        response: str, 
        redact_pii: bool = True
    ) -> Tuple[bool, str]:
        """
        Check and optionally redact output response.
        
        Returns:
            (is_safe, response_or_rejection_message)
        """
        import re
        
        modified_response = response
        found_pii = []
        
        # Detect and redact PII if requested
        if redact_pii:
            for pii_type, pattern in self.pii_patterns.items():
                matches = re.findall(pattern, modified_response, re.IGNORECASE)
                if matches:
                    found_pii.extend(matches)
                    modified_response = re.sub(
                        pattern, 
                        self.redaction_map[pii_type], 
                        modified_response, 
                        flags=re.IGNORECASE
                    )
        
        # Check for toxic content in response
        toxic_patterns = [
            "i hate", "stupid", "idiot", "worthless", "dumb",
            "kill", "hurt", "die", "threat", "attack"
        ]
        
        response_lower = modified_response.lower()
        for pattern in toxic_patterns:
            if pattern in response_lower:
                # Return safe fallback instead of toxic response
                return False, "I'm unable to provide an appropriate response. Please contact support."
        
        # If we redacted PII, log it but still return the redacted response
        if found_pii:
            print(f"PII redacted: {found_pii}")
        
        return True, modified_response
    
    def structured_output_guardrail(
        self, 
        response: Dict, 
        required_fields: List[str]
    ) -> Tuple[bool, Dict]:
        """
        Ensure structured output (JSON) has all required fields and valid types.
        
        Args:
            response: Dictionary from LLM
            required_fields: List of field names that must exist
            
        Returns:
            (is_valid, validated_or_corrected_response)
        """
        if not isinstance(response, dict):
            return False, {"error": "Response is not a dictionary"}
        
        # Check all required fields exist
        missing_fields = [f for f in required_fields if f not in response]
        if missing_fields:
            return False, {"error": f"Missing required fields: {missing_fields}"}
        
        # Validate confidence field if present
        if "confidence" in response:
            try:
                confidence = float(response["confidence"])
                response["confidence"] = min(max(confidence, 0.0), 1.0)  # Clamp to 0-1
            except (ValueError, TypeError):
                response["confidence"] = 0.5
        
        return True, response

# Complete pipeline with guardrails
class SafeLLMPipeline:
    """
    End-to-end safe LLM pipeline with guardrails at every step.
    """
    
    def __init__(self):
        self.guardrails = Guardrails()
        self.validator = DataValidator()
        self.responder = LLMResponder()
        self.hallucination_detector = HallucinationDetector()
    
    def process(
        self, 
        query: str, 
        chunks: List[str],
        conversation_history: List[Dict] = None
    ) -> Dict:
        """
        Complete safe processing pipeline.
        
        Steps:
        1. Input guardrail (block harmful queries)
        2. Validate retrieved chunks
        3. Get LLM response
        4. Output guardrail (redact PII, block toxic responses)
        5. Detect hallucination (async, non-blocking)
        """
        
        # STEP 1: Input Guardrail
        is_safe, result = self.guardrails.input_guardrail(query)
        if not is_safe:
            return {
                "success": False,
                "answer": "I cannot process this request.",
                "should_escalate": True,
                "reason": result
            }
        
        # STEP 2: Validate chunks (pre-LLM)
        should_escalate, reason = self.validator.should_escalate_to_human(
            query, chunks, 0.8
        )
        if should_escalate:
            return {
                "success": False,
                "answer": "I need to connect you to a human agent.",
                "should_escalate": True,
                "reason": reason
            }
        
        # STEP 3: Get LLM response
        response = self.responder.get_response_with_confidence(
            query, chunks, conversation_history
        )
        
        # STEP 4: Output Guardrail
        is_safe_output, safe_answer = self.guardrails.output_guardrail(
            response["answer"]
        )
        
        if not is_safe_output:
            return {
                "success": False,
                "answer": safe_answer,
                "should_escalate": True,
                "reason": "Output guardrail triggered"
            }
        
        response["answer"] = safe_answer
        
        # STEP 5: Hallucination detection (non-blocking, log only)
        # In production, run this async and log to database
        is_hallucinated, _, claims = self.hallucination_detector.check_hallucination(
            query, response["answer"], chunks
        )
        
        if is_hallucinated:
            # Log for monitoring, but don't block the response
            print(f"WARNING: Hallucination detected in response: {claims}")
            response["hallucination_detected"] = True
        
        response["success"] = True
        return response

# Example usage:
# pipeline = SafeLLMPipeline()
# result = pipeline.process(
#     query="How do I get a refund?",
#     chunks=["Refunds take 5-7 business days to process"],
#     conversation_history=[{"role": "user", "content": "I want my money back"}]
# )
# print(result["answer"])
# if result.get("should_escalate"):
#     print(f"Escalate: {result.get('reason')}")