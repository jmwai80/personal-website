#!/bin/bash
# Retries Oracle ARM VM provisioning every 5 minutes across all ADs.
# Stops and notifies on success. Logs all attempts.

eval "$(/opt/homebrew/bin/brew shellenv)"
export SUPPRESS_LABEL_WARNING=True

COMPARTMENT_ID="ocid1.tenancy.oc1..aaaaaaaal2rj3nmoeys5xx2665wryfe26qy6k527xn5yabxkhhxpexyqem7a"
SUBNET_ID="ocid1.subnet.oc1.iad.aaaaaaaay3pxqfccuginxiyp63hb3alwczaoozjswrvxqxh47kgttuy4xnma"
IMAGE_ID="ocid1.image.oc1.iad.aaaaaaaac6ozbxqea5kb7to5qu3asvnqj5f4j6gcxhxipeafefzpwtxm6mwa"
SSH_KEY="$HOME/.ssh/oracle_vm.pub"
LOG="$HOME/Projects/personal-website/scripts/oracle-retry.log"
SUCCESS_FILE="$HOME/Projects/personal-website/scripts/oracle-vm-details.txt"
ADS=("tuSy:US-ASHBURN-AD-1" "tuSy:US-ASHBURN-AD-2" "tuSy:US-ASHBURN-AD-3")

notify() {
  osascript -e "display notification \"$1\" with title \"Oracle VM\" sound name \"Glass\"" 2>/dev/null
}

log() {
  echo "[$(date '+%Y-%m-%d %H:%M:%S')] $1" | tee -a "$LOG"
}

log "=== Starting Oracle VM retry script ==="
notify "Oracle retry script started — will try every 5 minutes"

while true; do
  for AD in "${ADS[@]}"; do
    log "Trying $AD (4 OCPU / 24GB)..."

    RESULT=$(oci compute instance launch \
      --compartment-id "$COMPARTMENT_ID" \
      --availability-domain "$AD" \
      --shape "VM.Standard.A1.Flex" \
      --shape-config '{"ocpus":4,"memoryInGBs":24}' \
      --image-id "$IMAGE_ID" \
      --subnet-id "$SUBNET_ID" \
      --assign-public-ip true \
      --display-name "personal-website-vm" \
      --ssh-authorized-keys-file "$SSH_KEY" \
      --boot-volume-size-in-gbs 50 \
      2>&1)

    INSTANCE_ID=$(echo "$RESULT" | python3 -c "import sys,json; d=json.load(sys.stdin); print(d['data']['id'])" 2>/dev/null)

    if [ -n "$INSTANCE_ID" ]; then
      log "SUCCESS! Instance created: $INSTANCE_ID in $AD"
      log "Waiting for instance to get public IP..."

      # Wait up to 3 minutes for the IP
      for i in $(seq 1 18); do
        sleep 10
        PUBLIC_IP=$(oci compute instance list-vnics \
          --instance-id "$INSTANCE_ID" \
          --compartment-id "$COMPARTMENT_ID" \
          --query 'data[0]."public-ip"' --raw-output 2>/dev/null)

        if [ -n "$PUBLIC_IP" ] && [ "$PUBLIC_IP" != "null" ]; then
          log "Public IP: $PUBLIC_IP"
          cat > "$SUCCESS_FILE" <<EOF
VM provisioned successfully!
Instance ID: $INSTANCE_ID
Availability Domain: $AD
Public IP: $PUBLIC_IP
SSH command: ssh -i ~/.ssh/oracle_vm ubuntu@$PUBLIC_IP
Date: $(date)
EOF
          notify "VM is live! IP: $PUBLIC_IP — check oracle-vm-details.txt"
          log "Details saved to $SUCCESS_FILE"
          log "=== Script complete ==="
          exit 0
        fi
      done

      log "Instance created but IP not yet assigned. Check OCI console."
      notify "VM created! Check OCI console for IP — $INSTANCE_ID"
      exit 0
    else
      MSG=$(echo "$RESULT" | python3 -c "import sys,json; d=json.load(sys.stdin); print(d.get('message','unknown error'))" 2>/dev/null || echo "timeout/error")
      log "Failed on $AD: $MSG"
    fi
  done

  log "All ADs failed. Waiting 5 minutes before retry..."
  sleep 300
done
