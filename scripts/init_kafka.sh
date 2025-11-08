#!/bin/bash

set -e

# create kafka topic
echo -e 'Creating kafka topics'
kafka-topics --bootstrap-server ${KAFKA_HOSTNAME}:${KAFKA_DOCKER_PORT} --create --if-not-exists --topic ${KAFKA_SCRAPER_TOPIC} --replication-factor 1 --partitions 1

echo -e 'Successfully created the following topics:'
kafka-topics --bootstrap-server ${KAFKA_HOSTNAME}:${KAFKA_DOCKER_PORT} --list

# create test messages
echo -e 'Creating test messages in kafka topic:'
for (( i=10; i<=20; i++ ))
do
  echo "Creating message ${i}"

  echo '{"site_name":"vixar","msg":"MESSAGE:'"${i}"', Мы предлагаем долговечную и качественную продукцию по самым выгодным ценам. Подробная информация о продукции: 1. Сверхвысокая совместимость с чипами: T76 не только поддерживает распространенные на рынке чипы EEPROM, EPROM и FLASH, но также имеет надежную поддержку микроконтроллеров и программируемых логических устройств (GAL/CPLD) с возможностью настройки и добавления новых моделей чипов. Независимо от того, имеете ли вы дело с последовательными, параллельными ROM, SPI NAND или даже чипами EMMC/EMCP, T76 с легкостью с ними справляется.","date":"'"${i}"'-10-2025"}' | kafka-console-producer \
    --bootstrap-server "${KAFKA_HOSTNAME}:${KAFKA_DOCKER_PORT}" \
    --topic "${KAFKA_SCRAPER_TOPIC}"
done

echo -e 'Successfully created testing messages in kafka topic'