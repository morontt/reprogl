<?php

$fpIn = fopen(__DIR__ . '/logo.svg', 'r');
$fpOut = fopen(__DIR__ . '/logo__.svg', 'w');

while (($buffer = fgets($fpIn)) !== false) {
    echo $buffer;
    if (strpos($buffer, 'path') !== false) {
        $line = preg_replace_callback('/(\d+\.?\d*)/', function ($matches) {
            $in = (float)$matches[0];
            echo $in . "\n";
            $out = round($in * 512.0 / 310.0, 3);

            return $out;
        }, $buffer);
        fwrite($fpOut, $line);
    } else {
        fwrite($fpOut, $buffer);
    }
}

fclose($fpIn);
fclose($fpOut);
