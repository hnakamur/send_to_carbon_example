for i in sv{001..100}; do
  ./send_to_carbon_example -graphite-port 2003 -interval 1m -server-id $i -site-count 1000 &
done
wait
