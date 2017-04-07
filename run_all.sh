for i in sv{01..100}; do
  ./send_to_carbon_example -server-id $i -site-count 100 &
done
wait
