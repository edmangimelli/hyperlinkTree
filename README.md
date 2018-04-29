# hyperlinkTree
hyperlink a directory tree of html files recursively

## Description:
Given a directory, hyperlinkTree will walk through it, recursing through subdirectories, and create (or append to) "index.html", in each directory, with links to each html file and subdirectory in the current folder.
It will then, additionally, append the hyperlinks "previous", "up", and "next" to each html file (including the just-created index.html files). Folders that contain only a single html file and no subdirectories are treated specially, and do not get a generated index file—this is explained below (🐁).

## How it works:
### Step 1:
First hyperlinkTree will recurse through the directories and make a file named "index.html" in each directory.

For instance, if you had the following directories and files:
```
myDir/
   balloon.html
   car.html
   someFolder/
      lava.html
      88.HTM
      notHtml.txt
   tuna.Html
```
Two index files would be created:
```
myDir/
   balloon.html
   car.html
-> index.html
   someFolder/
      88.HTM
----> index.html
      lava.html
      notHtml.txt
   tuna.Html
```

The first index.html would be created with the following html (or this html would be appended to index.html, if index.html already existed):

```
myDir:<br>
<a href="balloon.html">balloon</a><br>
<a href="car.html">car</a><br>
<a href="someFolder/index.html">someFolder/</a><br>
<a href="tuna.html">tuna</a><br>
```
Simple relative links. Notice that the link "someFolder/" points to the next index.html.

The second index.html is similarly simple
```
someFolder:
<a href="88.HTM">88</a>
<a href="lava.html">lava</a>
```
Notice that the file "notHtml.txt" is not hyperlinked—hyperlinkTree is only interested in files with the extension "html" or "htm" (case-insensitive) or directories.


A more complicated example:
```
myBooks/
   Unfinished Book/
   Squirrelerella_Gets_Married/
      Ch1.html
      Ch2.html
      Ch3.html
      image01.jpg
      image02.jpg
      image03.jpg
   Writing_Performant_COBOL/
      Introduction.html
      Chapter_1/
         part1.html
         part2.html
         footNotes.html
      Chapter_2/
         text.HTM
         dog.jpg
         spaghetti.jpg
      Chapter_3.html
      Chapter_4/
         text.html
   Supporting Files/
      something.css
      scripty.sh
      moral.support
```
Here is what the index.html files would look like:

myBooks/index.html:
```
myBooks:<br>
<a href="Squirrelerella_Gets_Married/index.html">Squirrelerella_Gets_Married/</a><br>
<a href="Writing_Performant_COBOL/index.html">Writing_Performant_COBOL/</a><br>
```
Notice that two folders are unaccounted for: "Unfinished Book" and "Supporting Files". "Unfinished Book" is completely empty—so it was skipped. Likewise, "Supporting Files" has zero html files and zero subdirectories, making it appear empty to hyperlinkTree.

Squirrelerella_Gets_Married/index.html:
```
Squirrelerella_Gets_Married:<br>
<a href="Ch1.html">Ch1</a><br>
<a href="Ch2.html">Ch2</a><br>
<a href="Ch3.html">Ch3</a><br>
```
Nothing too interesting in that index.

Writing_Performant_COBOL/index.html
```
<a href="Introduction.html">Introduction</a><br>
<a href="Chapter_1/index.html">Chapter_1/</a><br>
<a href="Chapter_2/text.HTM">Chapter_2/text</a><br>
<a href="Chapter_3.html">Chapter_3</a><br>
<a href="Chapter_4/text.html">Chapter_4/text</a><br>

```
🐁 _The Special Rule:_ Notice that Chapter 2 and Chapter 4 both have links directly to text.HTM and text.html (respectively). When a subdirectory contains only one html file and no subdirectories, it is not given its own index.html file, rather, its single html file is linked in the parent's index.html. This functionality seemed convenient to me. It would be obnoxious to navigate to a folder only to have one possible option to select, so, those indices are foregone and the lonely html file is linked instead. This is especially convenient when you've organized your website so that each page is in its own subdirectory which contains the one html file for the page and all of its associated helper files (images, css, javascript, etc.).


Writing_Performant_COBOL/Chapter_1/index.html
```
Chapter_1:
<a href="part1.html">part1</a><br>
<a href="part2.html">part2</a><br>
<a href="footnotes.html">footnotes</a><br>
```
Nothing special in this index.

### Step 2:

After the indices have been created, the directory tree is recursively walked again, this time creating a list of all the html files in the tree. Continuing with the example above:

```
myBooks/
   index.html
   Unfinished Book/
   Squirrelerella_Gets_Married/
      index.html
      Ch1.html
      Ch2.html
      Ch3.html
      image01.jpg
      image02.jpg
      image03.jpg
   Writing_Performant_COBOL/
      index.html
      Introduction.html
      Chapter_1/
         index.html
         part1.html
         part2.html
         footNotes.html
      Chapter_2/
         text.HTM
         dog.jpg
         spaghetti.jpg
      Chapter_3.html
      Chapter_4/
         text.html
   Supporting Files/
      something.css
      scripty.sh
      moral.support
```

The list would be:
```
myBooks/index.html
Squirrelerella_Gets_Married/index.html
Squirrelerella_Gets_Married/Ch1.html
Squirrelerella_Gets_Married/Ch2.html
Squirrelerella_Gets_Married/Ch3.html
Writing_Performant_COBOL/index.html
Writing_Performant_COBOL/Introduction.html
Writing_Performant_COBOL/Chapter_1/index.html
Writing_Performant_COBOL/Chapter_1/part1.html
Writing_Performant_COBOL/Chapter_1/part2.html
Writing_Performant_COBOL/Chapter_1/footNotes.html
Writing_Performant_COBOL/Chapter_2/text.HTM
Writing_Performant_COBOL/Chapter_3.html
Writing_Performant_COBOL/Chapter_4/text.html
```
Notice that index.html is the first file for each directory on the list.

From here:
- `myBooks/index.html` would get appended a hyperlink with the text "next" that would point to `Squirrelerella_Gets_Married/index.html`.
- `Squirrelerella_Gets_Married/index.html` would get a "previous" link that pointed to `myBooks/index.html`, an "up" link that pointed again to `myBooks/index.html`, and a "next" link that pointed to `Squirrelerella_Gets_Married/Ch1.html`.
- ...
- `Writing_Performant_COBOL/Chapter_1/footnotes.html` would get a "previous" link to `Writing_Performant_COBOL/Chapter_1/part2.html`, an "up" link to `Writing_Performant_COBOL/Chapter_1/index.html`, and a "next" link to `Writing_Performant_COBOL/Chapter_2/text.HTM`
- and so on...

### Notes:
I realize that I'm hijacking a very commonly used file for websites—index.html—and this might be unideal for many real-world websites. The name "index.html" seemed most appropriate to me. This tool is first and foremost an educational exercise I dreamt up for myself. Cut me some slack—this is my first github project!

### To Do:
- At the moment you could put a dir in a dir in a dir in a dir ... and at the very end of chain is a non-html file, and this would make indices in all of the subdirectories up until the subdirectory with the subdirectory that has the non-html file—this is obviously undesirable.
