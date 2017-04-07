for i in sv{01..100}; do
  ./send_to_carbon_example -graphite-port 2005 -interval 1m -server-id $i -site-count 1000 &
done
wait
