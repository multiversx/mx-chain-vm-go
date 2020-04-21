#!/bin/bash

$ARWENDEBUG create-account --address=erdfoo --balance=100000 --nonce=42 || { exit 1; }
$ARWENDEBUG deploy --impersonated=erdfoo || { exit 1; }