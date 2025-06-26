#!/usr/bin/env python3
import requests
import sys

def test_health_endpoint(base_url):
    """Test the health endpoint returns 200 OK"""
    try:
        response = requests.get(f"{base_url}/health", timeout=5)
        assert response.status_code == 200, f"Expected 200, got {response.status_code}"
        assert response.text.strip() == "OK", f"Expected 'OK', got '{response.text.strip()}'"
        return True
    except Exception as e:
        print(f"Health check failed: {e}")
        return False

def test_root_endpoint(base_url):
    """Test the root endpoint returns expected message"""
    try:
        response = requests.get(base_url, timeout=5)
        assert response.status_code == 200, f"Expected 200, got {response.status_code}"
        assert "Hello World" in response.text, f"Expected 'Hello World' in response"
        return True
    except Exception as e:
        print(f"Root endpoint test failed: {e}")
        return False

if __name__ == "__main__":
    base_url = sys.argv[1] if len(sys.argv) > 1 else "http://localhost:8080"
    
    tests = [
        ("Health Endpoint", test_health_endpoint),
        ("Root Endpoint", test_root_endpoint)
    ]
    
    all_passed = True
    for test_name, test_func in tests:
        print(f"Running {test_name}...", end=" ")
        if test_func(base_url):
            print("✓ PASSED")
        else:
            print("✗ FAILED")
            all_passed = False
    
    sys.exit(0 if all_passed else 1)