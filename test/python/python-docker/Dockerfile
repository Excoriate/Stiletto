# First stage: Base Python environment
FROM python:3.9-slim AS base

ENV PYTHONDONTWRITEBYTECODE 1
ENV PYTHONUNBUFFERED 1

WORKDIR /app

# Second stage: Install dependencies
FROM base AS builder

RUN pip install --upgrade pip
COPY requirements.txt .
RUN pip install --user -r requirements.txt

# Third stage: Copy installed dependencies
FROM base AS runtime

COPY --from=builder /root/.local /root/.local

# Set the PATH to include local packages
ENV PATH=/root/.local/bin:$PATH

# Copy the FastAPI app
COPY main.py .

# Expose the FastAPI port
EXPOSE 8080

# Run the FastAPI app with Uvicorn
CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8080"]
