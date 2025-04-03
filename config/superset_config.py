from superset.security import SupersetSecurityManager

class CustomSecurityManager(SupersetSecurityManager):
    pass

# Отключаем телеметрию
FEATURE_FLAGS = {
    "ENABLE_TELEMETRY": False,
}

# Кеширование в Redis
CACHE_CONFIG = {
    'CACHE_TYPE': 'RedisCache',
    'CACHE_DEFAULT_TIMEOUT': 86400,
    'CACHE_KEY_PREFIX': 'superset_',
    'CACHE_REDIS_URL': 'redis://redis:6379/0'
}

# Асинхронные задачи через Celery
class CeleryConfig:
    broker_url = "redis://redis:6379/0"
    result_backend = "redis://redis:6379/0"
    accept_content = ["json"]
    task_serializer = "json"
    result_serializer = "json"

CELERY_CONFIG = CeleryConfig

# Безопасность
SESSION_COOKIE_SECURE = True
ENABLE_PROXY_FIX = True  # Если Superset за Nginx
