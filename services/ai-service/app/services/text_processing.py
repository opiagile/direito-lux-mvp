"""Text processing utilities for legal documents."""

import re
from typing import List, Dict, Any, Optional, Tuple
from unidecode import unidecode

from app.core.logging import get_logger

logger = get_logger(__name__)


class LegalTextProcessor:
    """Process and clean legal texts."""
    
    def __init__(self):
        """Initialize text processor."""
        # Legal reference patterns
        self.law_patterns = {
            "federal_law": re.compile(r"Lei\s+(?:Federal\s+)?n?º?\s*(\d+\.?\d*)/(\d{2,4})", re.IGNORECASE),
            "decree": re.compile(r"Decreto\s+(?:Federal\s+)?n?º?\s*(\d+\.?\d*)/(\d{2,4})", re.IGNORECASE),
            "constitution": re.compile(r"(?:CF|Constituição\s+Federal)(?:\s+de\s+\d{4})?", re.IGNORECASE),
            "article": re.compile(r"art(?:igo)?\.?\s*(\d+)", re.IGNORECASE),
            "paragraph": re.compile(r"§\s*(\d+)º?", re.IGNORECASE),
            "item": re.compile(r"inciso\s+([IVXLCDM]+|\d+)", re.IGNORECASE),
        }
        
        # Entity patterns
        self.entity_patterns = {
            "cpf": re.compile(r"\d{3}\.\d{3}\.\d{3}-\d{2}"),
            "cnpj": re.compile(r"\d{2}\.\d{3}\.\d{3}/\d{4}-\d{2}"),
            "process_number": re.compile(r"\d{7}-\d{2}\.\d{4}\.\d\.\d{2}\.\d{4}"),
            "monetary_value": re.compile(r"R\$\s*[\d.,]+", re.IGNORECASE),
            "date": re.compile(r"\d{1,2}/\d{1,2}/\d{2,4}"),
        }
        
        # Noise patterns to remove
        self.noise_patterns = [
            re.compile(r"^\s*Página\s+\d+\s*de\s+\d+\s*$", re.MULTILINE | re.IGNORECASE),
            re.compile(r"^\s*\d+\s*$", re.MULTILINE),  # Page numbers
            re.compile(r"_{3,}"),  # Long underscores
            re.compile(r"-{3,}"),  # Long dashes
            re.compile(r"\s{3,}"),  # Multiple spaces
        ]
    
    def clean_legal_text(self, text: str) -> str:
        """Clean and normalize legal text."""
        if not text:
            return ""
        
        # Remove common noise
        for pattern in self.noise_patterns:
            text = pattern.sub(" ", text)
        
        # Normalize whitespace
        text = " ".join(text.split())
        
        # Normalize legal citations
        text = self._normalize_legal_citations(text)
        
        # Remove excessive punctuation
        text = re.sub(r"([.!?]){2,}", r"\1", text)
        
        return text.strip()
    
    def _normalize_legal_citations(self, text: str) -> str:
        """Normalize legal citations for better matching."""
        # Normalize law references
        def normalize_law(match):
            number = match.group(1).replace(".", "")
            year = match.group(2)
            if len(year) == 2:
                year = "19" + year if int(year) > 50 else "20" + year
            return f"Lei {number}/{year}"
        
        text = self.law_patterns["federal_law"].sub(normalize_law, text)
        
        # Normalize article references
        text = re.sub(r"artigo\s+", "art. ", text, flags=re.IGNORECASE)
        
        return text
    
    def extract_legal_entities(self, text: str) -> Dict[str, List[str]]:
        """Extract legal entities from text."""
        entities = {
            "laws": [],
            "articles": [],
            "paragraphs": [],
            "process_numbers": [],
            "monetary_values": [],
            "dates": [],
            "cpf_cnpj": []
        }
        
        # Extract laws
        for match in self.law_patterns["federal_law"].finditer(text):
            law = f"Lei {match.group(1)}/{match.group(2)}"
            if law not in entities["laws"]:
                entities["laws"].append(law)
        
        # Extract articles
        for match in self.law_patterns["article"].finditer(text):
            article = f"Art. {match.group(1)}"
            if article not in entities["articles"]:
                entities["articles"].append(article)
        
        # Extract process numbers
        for match in self.entity_patterns["process_number"].finditer(text):
            if match.group() not in entities["process_numbers"]:
                entities["process_numbers"].append(match.group())
        
        # Extract monetary values
        for match in self.entity_patterns["monetary_value"].finditer(text):
            if match.group() not in entities["monetary_values"]:
                entities["monetary_values"].append(match.group())
        
        # Extract dates
        for match in self.entity_patterns["date"].finditer(text):
            if match.group() not in entities["dates"]:
                entities["dates"].append(match.group())
        
        # Extract CPF/CNPJ (anonymized)
        for pattern_name in ["cpf", "cnpj"]:
            for match in self.entity_patterns[pattern_name].finditer(text):
                anonymized = self._anonymize_document(match.group())
                if anonymized not in entities["cpf_cnpj"]:
                    entities["cpf_cnpj"].append(anonymized)
        
        return entities
    
    def _anonymize_document(self, doc: str) -> str:
        """Anonymize CPF or CNPJ."""
        if len(doc) == 14:  # CPF
            return doc[:3] + ".***.**" + doc[-3:]
        else:  # CNPJ
            return doc[:2] + ".***.***/****" + doc[-3:]
    
    def extract_key_phrases(self, text: str, max_phrases: int = 10) -> List[str]:
        """Extract key legal phrases from text."""
        # Simple implementation - can be enhanced with NLP
        sentences = re.split(r'[.!?]+', text)
        
        key_phrases = []
        for sentence in sentences:
            sentence = sentence.strip()
            
            # Look for sentences with legal terms
            if any(pattern.search(sentence) for pattern in self.law_patterns.values()):
                if len(sentence) > 20 and len(sentence) < 200:
                    key_phrases.append(sentence)
            
            if len(key_phrases) >= max_phrases:
                break
        
        return key_phrases
    
    def chunk_text(self, text: str, chunk_size: int = 1000, overlap: int = 200) -> List[str]:
        """Split text into overlapping chunks."""
        if not text or chunk_size <= 0:
            return []
        
        chunks = []
        start = 0
        text_length = len(text)
        
        while start < text_length:
            end = start + chunk_size
            
            # Try to break at sentence boundary
            if end < text_length:
                # Look for sentence end
                sentence_end = text.rfind(".", start, end)
                if sentence_end > start:
                    end = sentence_end + 1
            
            chunk = text[start:end].strip()
            if chunk:
                chunks.append(chunk)
            
            # Move start position with overlap
            start = end - overlap if end < text_length else text_length
        
        return chunks
    
    def highlight_matches(self, text: str, query: str, context_words: int = 10) -> List[str]:
        """Extract text snippets around query matches."""
        if not text or not query:
            return []
        
        # Normalize for matching
        text_lower = text.lower()
        query_lower = query.lower()
        query_words = query_lower.split()
        
        highlights = []
        
        # Find all occurrences
        for word in query_words:
            if len(word) < 3:  # Skip short words
                continue
            
            start = 0
            while True:
                pos = text_lower.find(word, start)
                if pos == -1:
                    break
                
                # Extract context
                words = text.split()
                char_count = 0
                word_start = 0
                
                for i, w in enumerate(words):
                    if char_count <= pos < char_count + len(w) + 1:
                        word_start = i
                        break
                    char_count += len(w) + 1
                
                # Get surrounding words
                context_start = max(0, word_start - context_words)
                context_end = min(len(words), word_start + context_words + 1)
                
                highlight = " ".join(words[context_start:context_end])
                if highlight not in highlights:
                    highlights.append(highlight)
                
                start = pos + 1
        
        return highlights[:5]  # Return top 5 highlights