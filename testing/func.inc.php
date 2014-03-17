<?php

/***

StockSloper PHP Functions
Jason Fowler
June 2013

***/

$prices = array('0.01 - 0.99','1.00 - 9.99','10.00 - 19.99','20.00 - unlimited');
$volumes = array('1000 - 49999','50000 - 249999','250000 - 999999','1000000 - unlimited');

# get a symbol

function getSymbol() {
        require_once('/home/stocksloper/db.inc');
        # smallish query - 11 stocks
        #$getsymbol= mysql_query("SELECT symbol FROM stock_symbol_name where symbol like 'STR%'") or trigger_error(mysql_error());
        # bigish query - 800 stocks
        #$getsymbol= mysql_query("SELECT symbol FROM stock_symbol_name where symbol like '%.TO'") or trigger_error(mysql_error());
        # this is the big momma
        $getsymbol= mysql_query("SELECT symbol FROM stock_symbol_name") or trigger_error(mysql_error());
                while($row = mysql_fetch_array($getsymbol)){
		$symbol = $row['symbol'];
		foreach($row AS $key => $value) { $row[$key] = stripslashes($value); }
				$data[] = $row;
			} //while
		 return $data;
mysql_free_result($getsymbol);

} // function getsymbol




function curl($url){
    $ch = curl_init();
    curl_setopt($ch, CURLOPT_URL,$url);
    curl_setopt($ch, CURLOPT_RETURNTRANSFER,1);
    return curl_exec($ch);
    curl_close ($ch);
} //function curl



function getStocks() {

require_once('/home/stocksloper/db.inc');

$today = date('Y-m-d');
$today_array = explode("-",$today);
$today_year = $today_array[0];
$today_month = $today_array[1];
$today_month = $today_month - 1;
$today_day = $today_array[2];
$today_day = $today_day - 1;

# Get a symbol

$stocks = getSymbol();

		foreach($stocks as $stock)
			{
	   		$symbol = $stock['symbol'];
	   		#echo $symbol."\n";

# Grab the data from Yahoo Finance

$close = "";
$volume = "";

        #dbg
	# dates work newer first then older.  d = newest month but make sure that jan = 0... eg, 5 24 = june 24, 8 15 = sept 15
        #$csv = "http://ichart.finance.yahoo.com/table.csv?s=$symbol&d=9&e=18&f=2013&g=d&a=06&b=01&c=2011&ignore=.csv";
	$csv = "http://ichart.finance.yahoo.com/table.csv?s=$symbol&a=$today_month&b=$today_day&c=$today_year";
        $yahoo = curl($csv);
        $data_array = explode("\n",$yahoo);
        array_shift($data_array);
        array_pop($data_array);
        foreach ($data_array as $data) {
		$yarr = explode(",",$data);
			if( isset($yarr[4]) && isset($yarr[5])) {
			$ydate = $yarr[0];
			#$open = $yarr[1];
			#$high = $yarr[2];
			#$low = $yarr[3];
			$close = $yarr[4];
			$volume = $yarr[5];
			#$adj_close = $yarr[6];
				} //if isset
	#dbg
	#echo $symbol . "," . $ydate .",". $close.",". $volume."\n";
		else {
		continue;
		}
		if ( ($close == "0") || ($volume == "0") )
		{
		continue;
		}
		else {
			#dbg
        		#echo $symbol . "," . $ydate .",". $close.",". $volume."\n";
               		# $log->lwrite("Added $symbol,$ydate,$close");
        #test
	   # $putstock= mysql_query("insert into stockhisttest (symbol,ydate,volume,close) values ('$symbol','$ydate','$volume',$close)") or trigger_error(mysql_error());
	# prod
                   $putstock= mysql_query("insert into stockhistory (symbol,ydate,volume,close) values ('$symbol','$ydate','$volume',$close)") or trigger_error(mysql_error());

			} //else putstock
                } // foreach data_array
        } //foreach stocks

} //function getstocks



function getLatest()	{

require_once('/home/stocksloper/db.inc');

$today = date('Y-m-d');

$deletedata = mysql_query("truncate table current") or trigger_error(mysql_error());

	$latest = mysql_query("insert into current (symbol, ydate, avgvol, close, slope, tradingdays, addate)
select distinct a.symbol, a.ydate, c.volume as avgvolume, a.close, round(b.slope,8) as slope, b.tradingdays, '$today'
from
(select latest.*
from stockhistory latest
inner join
  (select symbol, max(ydate) as maxdate
  from stockhistory
  group by symbol)
  grouplatest
  on latest.symbol = grouplatest.symbol
  and latest.ydate = grouplatest.maxdate) a
    join
    (select symbol, slope, tradingdays
    from slopedata) b
	join (select symbol, volume
	from avgvolume) c
    where a.symbol = b.symbol
	and a.symbol = c.symbol
	and c.volume > 25000
	and a.close <> 0
	and c.volume <> 0
	and b.slope <> 0
	and a.ydate between date_sub( curdate( ) ,interval 15 day ) and curdate( )
	") or trigger_error(mysql_error());



} //function getcurrent




# get the total number of days a stock has traded


function getDays ($symbol) {

require_once('/home/stocksloper/db.inc');

$result = mysql_query("select count(1) FROM stockhistory where symbol = '$symbol'")or trigger_error(mysql_error());
$row = mysql_fetch_array($result);

	$total = $row[0];
	return $total;
mysql_free_result($result);

} // function getdays


function deleteSlopeData() {

require_once('/home/stocksloper/db.inc');

# empty database to start fresh
        $deletedata = mysql_query("truncate table slopedata") or trigger_error(mysql_error());

}

# main meat of the thing - get the stock history and calculate the slope based on number of trading days (NTD)

function getRecursive($symbol,$ntd) {

$log = new Loggy();
$log->lfile('slope.log');

$tc = "";
$tdc = "";
$tdtd = "";
$ntd_1 = $ntd + 1;
$sumntd = $ntd * $ntd_1;
$sumntd = $sumntd / 2;

require_once('/home/stocksloper/db.inc');


#$getstockhistory = mysql_query("select @rownum:=@rownum+1 trading_day, sh.close
#                                from stockhistory sh,
#                                (select @rownum:=0) r
#                                where symbol = '$symbol'
#                                order by ydate desc limit $ntd") or trigger_error(mysql_error());

$getstockhistory = mysql_query("select t.trading_day, close from
(select @rownum:=@rownum+1 trading_day, sh.close, sh.ydate
                                from stockhistory sh,
                                (select @rownum:=0) r
                                where symbol = '$symbol'
                                order by ydate desc limit $ntd) t order by ydate asc") or trigger_error(mysql_error());


                        while($row = mysql_fetch_array($getstockhistory)){
                                 foreach($row AS $key => $value) { $row[$key] = stripslashes($value); }
                        #dbg
        #               echo $row['trading_day'] . " -  " . $row['close']."\n";
                        $tradingday = $row['trading_day'];
                        $closeprice = $row['close'];
                        $tdb = $tradingday * $closeprice;
                        $tdc = $tdb + $tdc;
                        $tdt = $tradingday * $tradingday;
                        $tdtd = $tdtd + $tdt;
                        $tc = $tc + $closeprice;
                        #dbg
                        #echo "TDC = $tdc, TDT = $tdt, TDTD = $tdtd\n";
                } //while getstockhistory
mysql_free_result($getstockhistory);
                        #dbg
                        #echo "Totals:  sumX = $sumntd, sumY = $tc, sumXY = $tdc, sumX2 = $tdtd\n";

	$preslope1a = bcmul($ntd,$tdc,8);
        $preslope1b = bcmul($sumntd,$tc,8);
        $preslope1 = bcsub($preslope1a,$preslope1b,8);
        $preslope2a = bcmul($ntd,$tdtd,8);
        $preslope2b = bcmul($sumntd,$sumntd,8);
        $preslope2 = bcsub($preslope2a,$preslope2b,8);
        $slope = bcdiv($preslope1,$preslope2,8);
	$slope = bcmul($slope,'-1',8);

		if (($slope >= -0.001) && ($slope <= 0.001)) {
		# need to reverse slope as per requirments
			$slope = bcmul($slope,'-1',8);
			#$log->lwrite("Slope $slope in range for $symbol with $ntd days");
	# put slope data into database
	$putslope= mysql_query("insert into slopedata (symbol,slope,tradingdays) values ('$symbol','$slope','$ntd')") or trigger_error(mysql_error());
                        #echo "slope $slope in range for $symbol\n";
                         } // if slope
                elseif ($ntd <=503) {
			$ntd++;
                     #   $log->lwrite("Slope NOT in range for $symbol with $ntd days");
			getRecursive($symbol,$ntd);
		} //elseif


} //function getRecursive


function getAvgVol() {

require_once('/home/stocksloper/db.inc');
	 $deletedata = mysql_query("truncate table avgvolume") or trigger_error(mysql_error());

# Get a symbol

$stocks = getSymbol();

                foreach($stocks as $stock)
                        {
                        $symbol = $stock['symbol'];
	#dbg
#	$symbol = 'AAPL';
                        #echo $symbol."\n";

# Grab the data from Yahoo Finance

$avgvolume = "";

        #dbg
	$csv = "http://download.finance.yahoo.com/d/quotes.csv?s=$symbol&f=a2";
        $volume = curl($csv);
	#print_r($volume);
        #dbg
        #echo $symbol . ",". $volume."\n";
            $putstock= mysql_query("insert into avgvolume (symbol,volume) values ('$symbol','$volume')") or trigger_error(mysql_error());

        } //foreach stocks

} //function getAvgVol


# Loggy Class for logging stuff

class Loggy{
    // define default log file
    private $log_file = 'stocksloper.log';
    // define file pointer
    private $fp = null;
    // set log file (path and name)
    public function lfile($path) {
        $this->log_file = $path;
    }
    // write message to the log file
    public function lwrite($message){
        // if file pointer doesn't exist, then open log file
        if (!$this->fp) $this->lopen();
        // define script name
        $script_name = pathinfo($_SERVER['PHP_SELF'], PATHINFO_FILENAME);
        // define current time
        $time = date('H:i:s');
        // write current time, script name and message to the log file
        fwrite($this->fp, "$time|$message\n");
    }
    // open log file
    private function lopen(){
        // define log file path and name
        $lfile = $this->log_file;
        // define the current date (it will be appended to the log file name)
        $today = date('Y-m-d');
        // open log file for writing only; place the file pointer at the end of the file
        // if the file does not exist, attempt to create it
        $this->fp = fopen($lfile . '_' . $today, 'a') or exit("Can't open $lfile!");
    }
}


?>