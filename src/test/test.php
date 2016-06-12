<?php
//$ip=1234;
#$bin=pack('V','1111');
#print_r(bin2hex($bin));
$abc=unpack('Vlen','0010');
print_r($abc);
//$ip='121.13.249.210';
//$ipdot = explode('.', $ip);
//$ip = pack('N', ip2long($ip));
//$ipdot[0] = (int) $ipdot[0];
//echo $ipdot[0];
//$ip = pack('N', "2030959058");
//print_r($ip);
//$hex = hex2bin("2030959058");
//print_r($hex);
//if($ip==$hex){
//	echo 'true';
//}