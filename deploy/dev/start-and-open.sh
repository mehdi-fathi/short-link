
# Wait for a few seconds to ensure services are up (adjust as needed)
sleep 3

# Open the browser to the application URL (replace with your URL)
if command -v xdg-open &> /dev/null
then
    xdg-open http://localhost:8080/index
elif command -v open &> /dev/null
then
    open http://localhost:8080/index
else
    echo "Please open http://localhost:8080/index in your browser"
fi