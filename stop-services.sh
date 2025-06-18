#!/bin/bash

# Script para parar todos os microserviços

echo "🛑 Parando microserviços Direito Lux..."

# Kill services using saved PIDs
if [ -f .auth.pid ]; then
    AUTH_PID=$(cat .auth.pid)
    if kill -0 $AUTH_PID 2>/dev/null; then
        echo "🔐 Stopping Auth Service (PID: $AUTH_PID)..."
        kill $AUTH_PID
    fi
    rm -f .auth.pid
fi

if [ -f .tenant.pid ]; then
    TENANT_PID=$(cat .tenant.pid)
    if kill -0 $TENANT_PID 2>/dev/null; then
        echo "🏢 Stopping Tenant Service (PID: $TENANT_PID)..."
        kill $TENANT_PID
    fi
    rm -f .tenant.pid
fi

if [ -f .process.pid ]; then
    PROCESS_PID=$(cat .process.pid)
    if kill -0 $PROCESS_PID 2>/dev/null; then
        echo "📋 Stopping Process Service (PID: $PROCESS_PID)..."
        kill $PROCESS_PID
    fi
    rm -f .process.pid
fi

if [ -f .datajud.pid ]; then
    DATAJUD_PID=$(cat .datajud.pid)
    if kill -0 $DATAJUD_PID 2>/dev/null; then
        echo "🔗 Stopping DataJud Service (PID: $DATAJUD_PID)..."
        kill $DATAJUD_PID
    fi
    rm -f .datajud.pid
fi

if [ -f .notification.pid ]; then
    NOTIFICATION_PID=$(cat .notification.pid)
    if kill -0 $NOTIFICATION_PID 2>/dev/null; then
        echo "📧 Stopping Notification Service (PID: $NOTIFICATION_PID)..."
        kill $NOTIFICATION_PID
    fi
    rm -f .notification.pid
fi

# Also kill any remaining processes by name (backup)
pkill -f "auth-server" 2>/dev/null || true
pkill -f "tenant-server" 2>/dev/null || true
pkill -f "process-server" 2>/dev/null || true
pkill -f "datajud-server" 2>/dev/null || true
pkill -f "notification-server" 2>/dev/null || true

echo "✅ All services stopped!"