#!/bin/bash

# Variables
TIMESTAMP=$(date +"%Y-%m-%d_%H-%M-%S")
SNAPSHOT_DIR="./new-snapshots"
SNAPSHOT_SCHEMA_FILENAME="${TIMESTAMP}_schema_snapshot.sql"
SNAPSHOT_DATA_AND_SCHEMA_FILENAME="${TIMESTAMP}_data_schema_snapshot.sql"
EMAIL_RECIPIENT="your-email@example.com"

# Ensure snapshot directory exists
mkdir -p $SNAPSHOT_DIR

# Run pg_dump inside the inventorydb container
docker exec inventorydb /bin/bash -c "pg_dump -U root -d inventorydb -s -f /tmp/${SNAPSHOT_SCHEMA_FILENAME}"
docker exec inventorydb /bin/bash -c "pg_dump -U root -d inventorydb -f /tmp/${SNAPSHOT_DATA_AND_SCHEMA_FILENAME} --create"

# Copy the snapshots from the container to the host
docker cp inventorydb:/tmp/${SNAPSHOT_SCHEMA_FILENAME} ${SNAPSHOT_DIR}/${SNAPSHOT_SCHEMA_FILENAME}
docker cp inventorydb:/tmp/${SNAPSHOT_DATA_AND_SCHEMA_FILENAME} ${SNAPSHOT_DIR}/${SNAPSHOT_DATA_AND_SCHEMA_FILENAME}

# Send email with attachments
# echo "Your database snapshot is attached." | mail -s "Database Snapshot ${TIMESTAMP}" -A ${SNAPSHOT_DIR}/${SNAPSHOT_SCHEMA_FILENAME} -A ${SNAPSHOT_DIR}/${SNAPSHOT_DATA_AND_SCHEMA_FILENAME} ${EMAIL_RECIPIENT}

# Clean up
# rm -f ${SNAPSHOT_DIR}/${SNAPSHOT_SCHEMA_FILENAME} ${SNAPSHOT_DIR}/${SNAPSHOT_DATA_AND_SCHEMA_FILENAME}
