## Autobackup manifests operator

# Сборка и установка
make build           # Собрать бинарник
make install         # Установить в систему

# Разработка
make dev             # Запуск с hot-reload (нужен air)
make test            # Запуск unit-тестов
make test-integration # Интеграционные тесты (требует k8s)

# Качество кода
make check           # Все проверки: fmt, vet, lint, security
make lint-fix        # Автоисправление линтинга

# Релиз
make release         # Сборка для всех платформ + checksums
make snapshot        # Создать development snapshot

# Очистка
make clean           # Очистка артефактов
make distclean       # Глубокая очистка (кэши)