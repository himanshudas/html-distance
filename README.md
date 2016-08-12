# html-distance

html-distance is a go library for computing the proximity of the HTML pages. The implementation similiarity fingerprint is Charikar's simhash. 

We used BK Tree (Burkhard and Keller) for verifying if a fingerprint is closed to a set of fingerprint within a defined proximity distance. 

Distance is the hamming distance of the fingerprints. Since fingerprint is of size 64 (inherited from hash/fnv), Similiarity is defined as 1 - d / 64.

In normal scenario, similarity > 95% (i.e. d>3) could be considered as duplicated html pages.


## Command Line Interface

```
Usage of html-distance:

    go run html-distance.go
```

Example
```
$ go run html-distance.go

Enter first url:
http://www.google.com
Enter second url:
http://www.facebook.com

Fetching http://www.google.com, Got 200

Fetching http://www.facebook.com, Got 200

Fingerprint1     100011001000101110010101111011000101011101100101111000000010
Fingerprint2  101101000011101100001010011010011000001010011101111100100100100

Feature Distance is 28.
Shingle factor is 2.
HTML Similarity is 56.25%

```
## Credits 

- The html-distance implementation is missing from the gryffin project:
 
```
  https://github.com/yahoo/gryffin
```
- This is proof of concept implementation for html-distance and could be used independent of gryffin project.

## Talks and Slides

- AppsecUSA 2015: [abstract](http://sched.co/3Vgm), [slide](http://go-talks.appspot.com/github.com/yukinying/talks/gryffin/gryffin.slide), [recording](https://youtu.be/IWiR2CPOHvc)


## Licence

Code licensed under the BSD-style license. See LICENSE file for terms.
