# requirements.txt
# openai==1.30.0
# langchain==0.1.0
# langchain-openai==0.0.5
# numpy==1.24.0
# rank_bm25  # for BM25

import os
from typing import List, Dict, Tuple
import numpy as np
from openai import OpenAI
from langchain_openai import OpenAIEmbeddings, ChatOpenAI
from langchain.prompts import ChatPromptTemplate
from langchain_core.output_parsers import PydanticOutputParser
from pydantic import BaseModel, Field
from rank_bm25 import BM25Okapi

# Set your API key
os.environ["OPENAI_API_KEY"] = "your-key-here"
client = OpenAI()