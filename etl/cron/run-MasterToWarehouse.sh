#!/bin/sh
export HOP_PROJECT_NAME=imunisasi
export HOP_PROJECT_FOLDER=/files/projects/imunisasi
export HOP_ENVIRONMENT_NAME=env
export HOP_ENVIRONMENT_CONFIG_FILE_NAME_PATHS=/files/projects/imunisasi/.env.json

hop-run \
  -p imunisasi \
  -e env \
  -j /files/projects/imunisasi/MasterToWarehouse/RunAll.hwf \
  -r local \
  -l Basic
