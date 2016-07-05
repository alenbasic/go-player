# go-player

<h3>Introduction</h3>

![alt tag](http://i.imgur.com/SjbwDFJ.gif)

<p>go-player is a basic web interface to <a href="http://elinux.org/Omxplayer">omxplayer</a> on the Raspberry Pi. 
You give it a filepath where you store all your video files and it will find anything ending in an extension put
in the extension list and then display all those files on a webpage for you to be able to search through and then 
click on to start the video. From there you can then send commands to the film (such as pausing, fast forwarding, 
rewinding and of course stopping the film completely).

<h3>Building and Installing</h3>
	
After you've downloaded it, cd into the directory and compile it with <code>go build</code> and run it with 
<code>./go-player</code>. Once it's running, go to <code>http://ip-of-RPi:8080</code> for the fun to begin. 
You can move the directory to wherever you want just make sure you include the HTML files with it as they're 
set up to look for a relative path (easily changeable of course in the code).</p>

<h3>To Do</h3>
<p>Some ideas I have in mind are (ordered somewhat in level of interest/priority):</p>
<ul>
<li><strike>The ability to add locations from within the web app itself. Right now it’s hard coded into the app 
before compilation. I was considering having a command line argument, but since I’ve sorta designed this to be run 
24/7 it makes more sense to be more flexible and have this looked after through the web interface. </strike></li>
<li>Related to this, it would probably be nice to clean up file names from within the interface. This is sorta 
breaking the UNIX philosophy here having an app do one thing and doing it well, and it would then require the app 
having more control over your system which makes me a bit more uncomfortable. This might evolve into an alias type 
system instead. We’ll see.</li>
<li>Being more robust in error handling. I’m currently having the app log more and more errors (where appropriate) 
to help troubleshoot issues as they happen, but it’d be nice to give this information through the web page itself 
rather than via the console (once again, where appropriate).</li>
<li>Writing up some unit (and other kinds of) tests. I think this is an important skill for me to work on and I 
believe that having a TDD mindset would be beneficial to me. Furthermore, it will save me manually having to test 
changes as I make them (which is a bit tedious as it involves me deploying it to my RPi, compiling it and then 
trying each of the pages out et cetera).</li>
<li>Better mobile support. Right now it looks a bit <em>meh</em> on anything other than a desktop/laptop which suits 
my needs but it really wouldn’t be a whole lot more effort to make it play ball with mobile.</li>
<li><strike>Create a GIF of the web page in action. I know I often look for previews of apps before I give them a go 
and will try an app that I’ve seen a preview of every time over one that’s only got documentation. This will involve 
some screen capture and editing. I’ll get to it at some point.</strike></li>
</ul>


<h3>Background</h3>

<p>This is the go version of my other project, <a href="https://github.com/alenbasic/sto-player">sto-player</a> 
which is written in python and uses web.py and popen as external dependencies. One of the reasons I started this 
project was I noticed a bit of slowness to sto-player on occasion (Testing this on the RPi2 hardware). sto-player 
could sometimes take around 2 seconds to perform commands like starting or stopping a film. go-player in contrast 
however is pretty much instantaneous in everything it does. </p>

<p>Note: You may think that this is due to the fact that I saved the file information into a sqlite database in 
sto-player, but this is even when we take this into account by removing that functionality</p>

<h3>Requirements</h3>

go-player depends has no external dependencies and uses only the go standard library, so if you have go installed, 
that's all you need to compile the program and get started.

<h3>How to Use</h3>

As mention in the intro, you give it your media directory location, navigate to <code>http://ip-of-RPi:8080</code> in 
your web browser and select a movie from there. The interface is pretty straight forward I believe, so I don't think there 
is much of a need of an explanation beyond that. Since there is always room for improvement, feel free to contact me 
regarding improvements or suggestions to the interface or the application itself.

<h3>License</h3>
go-player is released under the v3 of the GPL, please refer to the LICENSE file for more information.