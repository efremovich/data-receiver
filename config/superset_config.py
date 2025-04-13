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

# Добавление русского языка в список доступных языков
LANGUAGES = {
    "ru": {"flag": "ru", "name": "Русский"},
    "en": {"flag": "us", "name": "English"}
}
# Установка русского языка в качестве языка по умолчанию
BABEL_DEFAULT_LOCALE = "ru"

# Установка таймаута до 5 минут
SQLALCHEMY_QUERY_TIMEOUT = 300
SUPERSET_WEBSERVER_TIMEOUT = 600

FEATURE_FLAGS = {
    "ALERT_REPORTS": True,
    "DATAPANEL_CLOSED_BY_DEFAULT": True,
    "DASHBOARD_VIRTUALIZATION": True,
    "DASHBOARD_RBAC": True,
    "ENABLE_TEMPLATE_PROCESSING": True,
    "ESCAPE_MARKDOWN_HTML": True,
    "LISTVIEWS_DEFAULT_CARD_VIEW": True,
    "THUMBNAILS": True,
    "DRILL_BY": True,
    "DRILL_TO_DETAIL": True,
    "HORIZONTAL_FILTER_BAR": True,
    "ESTIMATE_QUERY_COST": True,
    "TAGGING_SYSTEM": True,
    "HTML_SANITIZATION": False,
}
