#!/bin/bash

IFS=',' # set the field separator to comma

i=0
tail -n +2 wallets.csv | while IFS=',' read -r address key
do
  # shellcheck disable=SC2059
  printf "Deploy contract for index $i\n"
  sed "s/{{WALLET_KEY}}/$key/g" ./hardhat-boilerplate/.env.tmp > ./hardhat-boilerplate/.env
  rand_str=$(cat /dev/urandom | LC_ALL=C tr -dc 'A-Z' | fold -w 3 | head -n 1)
  sed "s/{{TOKEN_NAME}}/$rand_str/g" ./hardhat-boilerplate/contracts/Token.sol.tmp > ./hardhat-boilerplate/contracts/Token.sol

  # shellcheck disable=SC2164
  cd hardhat-boilerplate
  npx hardhat compile
  cli_output=$(npx hardhat run scripts/deploy.ts --network base-goerli)
  contract_address=$(echo "$cli_output" | grep -oE '0x[a-fA-F0-9]{40}' | cut -d ' ' -f 1)
  # shellcheck disable=SC2103
  cd ..
  echo "$address" >> blacklist.csv
  echo "$address, $key, $contract_address" >> output.csv
  printf "\n"
  i=$((i+1))
done

printf "Deploy contract successfully"