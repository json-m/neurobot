#!/bin/zsh
response=$(curl -s 'http://100.71.7.20:9292/cs' | jq '.')
hits=$(echo "${response}" | jq '.hit')
misses=$(echo "${response}" | jq '.miss')
lookups=$(echo "${response}" | jq '.lookups')
if (( hits + misses > 0 )); then
    hit_ratio=$(bc <<< "scale=2; $hits * 100 / $lookups")
    echo "Cache hit ratio is $hit_ratio% (H:$hits|M:$misses|T:$lookups)"
else
    echo "No cache hits or misses recorded yet"
fi
