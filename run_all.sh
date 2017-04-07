for i in sv{01..49}; do
  ./send_to_carbon_example -server-id $i &
done
wait
