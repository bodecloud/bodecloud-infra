#!/bin/bash

echo "Finding pods with image pull errors..."
kubectl get pods -A | grep -E "ImagePullBackOff|ErrImagePull" | awk '{print $2}' > failing_pods.txt

echo "Extracting bad images from failing pods..."
echo "# Bad images from failing pods - $(date)" > bad_images.txt

while read pod; do
  echo "Processing pod $pod..."
  kubectl get pod $pod -o=jsonpath='{range .spec.containers[*]}{.image}{"\n"}{end}' 2>/dev/null >> bad_images.txt
done < failing_pods.txt

echo "Removing duplicates..."
sort bad_images.txt | uniq > temp.txt && mv temp.txt bad_images.txt

echo "Found $(grep -v "#" bad_images.txt | wc -l) bad images. Results saved to bad_images.txt" 