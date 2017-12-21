# gget
A specialized replacement of wget in Go

## To Build
`go build -o gget src/cmd/gget/main.go`

## Background
I had a special application that required the fetching of a lot of files listed in a text file.
Normally in cases like this I would use wget -i, but since the urls contained unicode characters needing URL-encoding, 
Wget messed up when naming files and folders by failing to decode it back for a proper filename.

So here is my very go-ish solution to the problem.

## Workflow
A file where each line is a URL to fetch is loaded and each line is put in a work-queue to processing.
a N amount of workers consumes the queue concurrently, fetching files over HTTP and storing them in the filesystem.

## Input
Example input:

````
http://sigmaitc.se/wp-content/uploads/sites/4/2014/02/Panelbild_jobs-800x333.jpg
http://sigmaitc.se/wp-content/uploads/sites/4/2015/03/sigma-compleo.jpg
http://sigmaitc.se/wp-content/uploads/sites/4/2014/10/mailbanner-kille.jpg
http://sigmaitc.se/wp-content/uploads/sites/4/2014/02/medarbetare_tomas_hellman.jpg
http://sigmaitc.se/wp-content/uploads/sites/4/2014/02/medarbetare_lars_littorin.jpg

http://www.kammarkollegiet.se/sites/default/files/2014-06-12, Remissyttrande - Översvämningsmyggor vid Nedre Dalälven (PDF, 1.01 MB).pdf
````

## Output
Nothing if everything went well. A list of errors if something went wrong.

### Side effects
For each domain, a folder structure is created that mimics the structure of the URLs.
So given the example input above, the following files would be created (Windows path assumed, but should support Unix as well on relevant systems.)
````
sigmaitc.se\wp-content\uploads\sites\4\2014\02\Panelbild_jobs-800x333.jpg
sigmaitc.se\wp-content\uploads\sites\4\2014\02\medarbetare_lars_littorin.jpg
sigmaitc.se\wp-content\uploads\sites\4\2014\02\edarbetare_tomas_hellman.jpg
sigmaitc.se\wp-content\uploads\sites\4\2014\10\mailbanner-kille.jpg
sigmaitc.se\wp-content\uploads\sites\4\2015\03\sigma-compleo.jpg
www.kammarkollegiet.se\sites\default\files\2014-06-12, Remissyttrande - Översvämningsmyggor vid Nedre Dalälven (PDF, 1.01 MB).pdf
````


