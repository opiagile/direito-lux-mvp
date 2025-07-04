#!/usr/bin/env python3

import json
import time
from http.server import HTTPServer, BaseHTTPRequestHandler
from urllib.parse import urlparse, parse_qs

class ProcessStatsHandler(BaseHTTPRequestHandler):
    def do_OPTIONS(self):
        self.send_cors_headers()
        self.end_headers()
    
    def do_GET(self):
        parsed_url = urlparse(self.path)
        
        if parsed_url.path == '/health':
            self.handle_health()
        elif parsed_url.path == '/api/v1/processes/stats':
            self.handle_stats()
        else:
            self.send_error(404, 'Not Found')
    
    def send_cors_headers(self):
        self.send_response(200)
        self.send_header('Access-Control-Allow-Origin', '*')
        self.send_header('Access-Control-Allow-Methods', 'GET, POST, PUT, DELETE, OPTIONS')
        self.send_header('Access-Control-Allow-Headers', 'Origin, Content-Type, Accept, Authorization, X-Tenant-ID')
        self.send_header('Content-Type', 'application/json')
    
    def handle_health(self):
        self.send_cors_headers()
        self.end_headers()
        
        response = {
            "status": "healthy",
            "service": "process-service-temp",
            "timestamp": time.time()
        }
        
        self.wfile.write(json.dumps(response).encode())
    
    def handle_stats(self):
        tenant_id = self.headers.get('X-Tenant-ID')
        if not tenant_id:
            self.send_error(400, 'X-Tenant-ID header é obrigatório')
            return
        
        # Dados diferentes por tenant - CAMPOS COMPLETOS PARA O DASHBOARD
        tenant_stats = {
            "11111111-1111-1111-1111-111111111111": {
                "total": 45, 
                "active": 38, 
                "paused": 5, 
                "archived": 2, 
                "this_month": 12,
                "todayMovements": 3,
                "upcomingDeadlines": 7
            },
            "22222222-2222-2222-2222-222222222222": {
                "total": 32, 
                "active": 28, 
                "paused": 3, 
                "archived": 1, 
                "this_month": 8,
                "todayMovements": 2,
                "upcomingDeadlines": 4
            },
            "33333333-3333-3333-3333-333333333333": {
                "total": 67, 
                "active": 58, 
                "paused": 7, 
                "archived": 2, 
                "this_month": 15,
                "todayMovements": 5,
                "upcomingDeadlines": 9
            },
            "44444444-4444-4444-4444-444444444444": {
                "total": 29, 
                "active": 25, 
                "paused": 2, 
                "archived": 2, 
                "this_month": 7,
                "todayMovements": 1,
                "upcomingDeadlines": 3
            },
            "55555555-5555-5555-5555-555555555555": {
                "total": 18, 
                "active": 16, 
                "paused": 1, 
                "archived": 1, 
                "this_month": 4,
                "todayMovements": 0,
                "upcomingDeadlines": 2
            },
            "66666666-6666-6666-6666-666666666666": {
                "total": 89, 
                "active": 76, 
                "paused": 8, 
                "archived": 5, 
                "this_month": 22,
                "todayMovements": 6,
                "upcomingDeadlines": 12
            },
            "77777777-7777-7777-7777-777777777777": {
                "total": 156, 
                "active": 142, 
                "paused": 10, 
                "archived": 4, 
                "this_month": 35,
                "todayMovements": 8,
                "upcomingDeadlines": 18
            },
            "88888888-8888-8888-8888-888888888888": {
                "total": 78, 
                "active": 71, 
                "paused": 5, 
                "archived": 2, 
                "this_month": 18,
                "todayMovements": 4,
                "upcomingDeadlines": 6
            }
        }
        
        stats = tenant_stats.get(tenant_id, {
            "total": 25, 
            "active": 20, 
            "paused": 3, 
            "archived": 2, 
            "this_month": 6,
            "todayMovements": 1,
            "upcomingDeadlines": 2
        })
        
        print(f"📊 Stats para tenant {tenant_id}: {stats}")
        
        self.send_cors_headers()
        self.end_headers()
        
        response = {
            "data": stats,
            "timestamp": time.time()
        }
        
        self.wfile.write(json.dumps(response).encode())

if __name__ == '__main__':
    port = 8083
    server = HTTPServer(('localhost', port), ProcessStatsHandler)
    print(f"🚀 Process Stats Server rodando na porta {port}")
    print(f"📊 Endpoint: GET /api/v1/processes/stats")
    print(f"🔍 Health: GET /health")
    
    try:
        server.serve_forever()
    except KeyboardInterrupt:
        print("\n⏹️ Servidor parado")
        server.shutdown()