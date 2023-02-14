go build

if [[ $? -eq 0 ]]
then
    ./openai-tg-bot $@
fi

