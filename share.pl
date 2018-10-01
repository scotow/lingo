#!/usr/bin/perl

use File::Basename;

sub usage {
	print 'usage: ', basename($0), " [NAME [URL]]\n";
}

if ($ARGV[0] eq '-h' || $ARGV[0] eq '--help') {
	usage;
	exit 0;
}

my $name = $ARGV[0];
my $url;

if ($ARGV[1] eq '') {
	$url = <STDIN>;
	chomp $url;
} else {
	$url = $ARGV[1];
}

exec 'curl', '-L', '-d', "url=$url", "https://share.scotow.com/$name";
