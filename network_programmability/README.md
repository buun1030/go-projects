# Data Formats and Models
If you’ve done any amount of exploration into the world of APIs, you’ve likely heard about data formats like JSON, XML, or YAML. You may have heard about concepts like data modeling, or model-driven APIs. Terms like data serialization and markup language may have popped into the foreground. You’d be right to wonder what all of this means and how it all applies to network automation.

It turns out that these concepts are at the heart of any reasonably complex modern software system, including those built and operated for the purpose of network automation. Even if you’re writing a simple script to change the hostname on a switch, at some point, your script will need to transmit some kind of information over the network that the switch will successfully receive and correctly interpret. How can you get your script and that switch to speak the same language?

Data formats like the aforementioned are those shared languages. They are broadly supported in all popular programming languages and are under the covers of nearly all the libraries and tools that you’ll use in your network automation journey. They are used by your network device’s built-in software for the purpose of being able to reliably and programmatically communicate with external entities, whether a full-blown fabric manager or a simple script on your laptop.

Understanding these formats, and how to work with the data they represent, is therefore crucial for you to be able to work effectively as a network automation professional. This chapter covers a variety of technologies and tools used to represent, transmit, store, and model data formats, with a specific focus on those that you’re most likely to run into in your network automation work.

## Benefits and Fundamentals of Structured Data Formats
A programmer typically uses a wide variety of tools to store and work with data. You could use simple scalar values (single values), collections (arrays, hashmaps), or even custom types built in the syntax of the language you’re using. While the specifics often differ, all languages offer primitives like this to give the programmer multiple ways to solve problems. When passing data within the context of a single program, this is often sufficient. A compiler knows exactly how much memory to allocate for a given type, so all you have to do as the programmer is reference that type when you need it, and the compiler will handle everything.

However, sometimes a more abstract, portable format is required. For instance, a non-programmer may need to be able to feed data into, or retrieve data from, a running program. Multiple programs may need to communicate with one another somehow—and the programs may not even be written in the same language; this is often the case with traditional client-server applications using a script you’ve written to automate a task on a network device, for example.

The data formats discussed in this chapter were designed to enable these kinds of use cases. They are well-established standards for communicating between generic software systems, and as a result, they’re well supported in any language or tool you choose to use. They give you the ability to describe data that would otherwise be represented as a series of bytes in memory.

NOTE
Without standardization of data formats, our networks wouldn’t even function! Protocols like BGP, OSPF, and TCP/IP were standardized out of a necessity for network elements to have a predictable, shared language in order to effectively communicate across a globally distributed system—​the internet!

The formats discussed in this chapter have three key traits that make them extremely useful and preferable, especially within the context of network automation:

Structured
These data formats, based on an agreed-upon set of rules, were designed to be easier for machines to understand. Computers are much more literal than humans and can’t intuitively understand data without a strict, predictable structure. For instance, the unstructured data you might see in the output of a show command on your router or switch may be formatted well for human consumption, but is not ideal for a computer to readily parse and understand.

Supported
Since these formats are standardized and widely adopted, you’ll almost never have to write your own code to understand them directly. You can reuse existing (and often extremely mature) software and tools for this. Many programming languages like Python and Go have built-in mechanisms that make it easy to import and export data to these formats, either on the filesystem or on the network.

Portable
While some languages have their own intermediate representations (i.e., pickle in Python or gobs in Go), the formats we discuss here are language agnostic, meaning they work with a wide variety of software ecosystems.

These are all important to consider, but let’s ponder the first point a little longer through some examples. Why are structured data formats easier for computers to understand, as opposed to the output you might see as a result of a simple show command?

Whenever you run a command like this, the software on your network device first gathers any data it needs from its subsystems or other network devices. At this stage, the information is little more than bytes in memory. To meaningfully display the results, the software then represents that information in a format that a human being can easily and quickly understand:

root@vqfx1> show interfaces em0
Physical interface: em0    , Enabled, Physical link is Up
  Interface index: 8, SNMP ifIndex: 17
  Type: Ethernet, Link-level type: Ethernet, MTU: 1514, Speed: 1000mbps
  Device flags   : Present Running
  Interface flags: SNMP-Traps
  Link type      : Full-Duplex
  Current address: 52:54:00:b1:f5:8d, Hardware address: 52:54:00:b1:f5:8d
  Last flapped   : 2019-01-10 17:49:55 UTC (00:17:33 ago)
    Input packets : 1039
    Output packets: 778
The nice thing about output like this is that it requires little effort (or even expertise) to see that the name of this interface is em0. Our brains have the tools to flexibly identify data by using helpful phrases like Physical interface. Even if what we have in mind isn’t exactly this, we know that it will get us what we need.

However, it’s actually not obvious to a computer where the interface name is located in this output. To us, the term Physical interface is a useful indicator to describe the nature of the text that follows. To a computer, it’s all just undifferentiated text. If you were to write a program to pull out the bits of valuable data from this output, you’d have to answer some important questions:

How do you know which portion of the text represents the value you want to access? Is it before or after the colon? What about commas? Why do some values share a line, whereas others get their own dedicated line?

What happens if the output doesn’t follow a consistent set of formatting rules?

What happens when another command (e.g., show bgp neighbor) formats things differently? Do you have to write a separate program or function for each command?

Since the primary use case for this output is human readability, what happens when the network vendor hires a UX expert to review and make changes to the format of this output?

When writing your parsing program or function, you’ll have to answer these questions, and more often than not, that will require a lot of extra time and energy that you may not be able to afford. In contrast, structured formats like JSON and XML were built to handle these concerns well.

For example, some configuration models are friendly to automated methods, by representing the configuration model in these data formats like XML or JSON. It is easy in Junos OS to see the XML representation of the show command we ran earlier, as shown in Example 8-1.

NOTE
Full versions of the code examples in this chapter can be found in the book’s GitHub repo at https://github.com/oreilly-npa-book/examples/tree/v2/ch08-dataformats.

Example 8-1. Displaying the XML-RPC equivalent for Junos commands
root@vqfx1> show interfaces em0 | display xml
<rpc-reply xmlns:junos="http://xml.juniper.net/junos/15.1X53/junos">
    <interface-information>
        <physical-interface>
            <name>em0</name>
            <admin-status junos:format="Enabled">up</admin-status>
            <oper-status>up</oper-status>
            <local-index>8</local-index>
            <snmp-index>17</snmp-index>
            <if-type>Ethernet</if-type>
            <link-level-type>Ethernet</link-level-type>
            <mtu>1514</mtu>
            <speed>1000mbps</speed>

 ... output truncated for brevity ...
Comparing this output to the preceding example, you might point out that this is quite a bit harder to read, and you’d be right. You’d even be right to argue that the latter is potentially a little less efficient; in some cases, more raw text would be needed to represent the same data.

From a programmatic perspective however, this is ideal. Example 8-1 provides key advantages over the previous, human-readable code when it comes to being able to programmatically parse the data contained within:

XML follows a stable set of rules, so there’s no need to constantly rewrite low-level text-parsing logic. You need to care about only the data being represented.

There is consistent use of delimiting structures. It’s clear that any tag starts with < and ends with >, and an opening tag should eventually be closed with a corresponding closing tag, such as </tag>.

Each piece of data is given its own easily parsable field. You know that the space inside the tags represents the entirety of the actual value, and everything else is just structure.

This format is inherently hierarchical. You know based on the order of opening and closing tags which values have parent/child relationships.

There is an established convention for metadata—that is, data about data (e.g., the xmlns tag).

While some of the particulars here are specific to XML, all structured data formats provide the same advantages in their own way.

So, in short, structured data formats like those discussed here are designed to allow software systems to communicate reliably and predictably with one another, no matter what language they’re written in.

### When Structured Data Isn’t Available: Screen Scraping
When the first version of this book was written, numerous platforms offered only human-readable text as output to be consumed by automation tools and scripts, as opposed to structured formats like JSON or XML. In cases like these, screen scraping can be used to retrieve data from a network device. This technique uses a protocol like SSH to emulate user behavior by sending a series of terminal commands, retrieving the raw text output, and attempting to format this output into a more structured representation. Whether done using a language like Python, or other tools that may provide a slightly more abstract framework, this approach requires you to provide your own set of low-level rules for parsing raw text, or depend on those created/maintained by others. However, these days, platforms that require this approach are becoming more and more rare.

Screen scraping not only makes your automation software extremely fragile, but is also enormously wasteful to your time, and by extension, that of the organization you’re working in. The lesson to learn here isn’t that it is impossible to write a program to parse some kind of text blob that follows from a show command. Instead, remember that the unstructured output you see in your terminal—and the subsystems required to produce it—simply weren’t designed to be accessed programmatically. Even those who have successfully created screen-scraping scripts will tell you that this approach is fraught with danger. Some network platforms will literally crash if too many commands are sent at once—not exactly a solid foundation for automation.

In contrast, choosing an architecture that aligns best with your own use case means you can achieve a valuable outcome more quickly and avoid having to deal with problems that were solved decades ago. Your job as network automation professionals is to provide value to your organization as quickly and effectively as possible, and if you are able to choose an approach that doesn’t require you to constantly reinvent parsing logic, you should.

Platforms that don’t offer any form of support for structured data formats are being phased out, and any serious automation initiative should include a requirement for platforms that do support these options. This not only saves you the time from having to do screen scraping yourself, but also frees you from having to use tools built on this fragile foundation.

### Types of Data
As discussed in Chapters 6 and 7, you can use a variety of built-in data types in any modern programming language. We refer to these throughout the text, so if you’re not familiar with them, we recommend you start there. We might use terms like string, integer, and boolean for representing different types of scalar (singular) values, list or array to describe a collection of values, or dictionary for key-value pairs.

However, all of these may be known by slightly different names in the various data formats and programming languages we reference, and within those contexts we may use different terms. This is OK and expected; it’s more important that you understand the basic concepts behind all of these, rather than trying to be unified and precise in the terminology across the board.

### Documents Versus Data
You may have heard the term markup language within the context of some of these formats. This is an important term to understand because it is a big part of the history of some of the formats we discuss, and we should be clear about the primary reasons we’re even talking about these formats in the first place.

Markup languages can also be referred to as document-oriented languages. The canonical example for this is HTML, which includes tags for things like headers, images, and links to external dependencies like JavaScript and CSS files:

<html>
    <body>
        <div>
            <p>Hello, World!</p>
        </div>
    </body>
</html>
Markup languages like HTML are ultimately used to describe/annotate a document, which is then rendered together to form a web page that shows up in your browser.

However, this chapter is not focused on this use case. While some formats (in particular XML) can be used for this purpose, the data formats we’re looking at are designed for the task of data serialization—that is, representing data (not documents) in a structured way. Remember that the primary reason we’re looking at these formats is so that you can understand how software systems exchange data with one another.

So in short, you can think of markup languages as describing documents and data serialization formats as describing data. It is this second use case that we’re focusing on in this chapter.

### Categories of Data Formats
The data formats we explore in this chapter fall into two broad categories:

Text-based
Data is first serialized into an intermediate format like UTF-8 and then encoded into bytes for storage or transmission.

Binary
Data is encoded directly into an efficient, binary format.

We’ll start by looking at text-based data formats.

## Text-Based Data Formats
We’ve already teased a few text-based formats so far in this chapter, but now it’s time to look at them more closely. Text-based formats have some key advantages:

You can easily edit them using a standard text editor or view them plainly using inspection tools in your browser.

They are well-established standards, and it’s extremely easy to find support for them in libraries and tools.

They’re abstract enough to map into just about any common data structure in a variety of programming languages.

NOTE
The main disadvantage to these formats is that they can be inefficient. One reason is because text-based formats include not only the raw data you want to transmit (strings, integers, arrays, key-value pairs) but also the various characters used to represent that data, such as curly braces and square brackets in JSON or <> tags in XML.

Text-based formats generally have to use more raw storage or bandwidth capacity to accommodate this extra information. These formats also require more processing to both send and receive data. This can become a problem when sending large amounts of data, as this inefficiency can compound. Fortunately, for the vast majority of network automation use cases and workflows, this is rarely a problem.

Whether at rest or in transit, any data you bring into an automation tool or script is ultimately represented as 0s and 1s: bits. Most of the time, you deal with these in multiples of 8: bytes. We loosely refer to this as raw binary data. However, before you can do something practical with this data, it needs to be processed and converted into a form you can work with. For example, computers don’t implicitly know that a series of bytes you’re receiving from an API request is ultimately meant to be interpreted as a Python list.

To get data from this raw binary format to something you can use (or vice versa), two distinct phases must take place:

Decoding and deserializing

Serializing and encoding

Figure 8-1 shows that to make sense of data using one of these formats—say, as a payload in an API response that you’re receiving—your computer must first decode the raw bits that come off the wire and into a text-encoding standard like UTF-8. This is the decoding step.

npa2 0801
Figure 8-1. Decoding and deserializing text-based data formats
However, at this point you effectively have the equivalent of one big string. A JSON payload, for example, might look something like Example 8-2.

Example 8-2. Decoded JSON string
"{\"vendors\":[\"Cisco\",\"Juniper\",\"Arista\"]}"
Data that’s been transmitted using these formats is most useful when it’s been deserialized into types and structures within the programming language or tool you’re using. For instance, this JSON document would map nicely into a Python dictionary, with a single key, vendors. This key’s value would map into a Python list for the three elements in the JSON array. Only after the data is decoded and then deserialized can you do something useful with the data contained within the response payload.

The same process must be followed in reverse to store or transmit data from these types. First, the data must be serialized into one of these formats and then encoded into bytes, as in Figure 8-2.

npa2 0802
Figure 8-2. Serializing and encoding of text-based data formats
This approach sacrifices a bit of efficiency for enhanced portability, and the potential for humans to more easily understand and even make changes to data found in one of these intermediate formats. Next, we’ll explore some of these specific formats in greater detail.

YAML
If you’re reading this book because you’ve seen compelling examples of network automation online or in a presentation and you want to learn more, you may have heard of YAML. This is because YAML is a particularly human-friendly data format, and for this reason, it is used in many network automation tools and initiatives. For instance, YAML is used by Ansible to describe playbooks, variable files, inventory files, and more, as you’ll see in Chapter 12.

NOTE
Previously, we explored the difference between markup and data serialization formats, and you may be wondering which category best describes YAML. Fortunately, the website that hosts the YAML specification (https://www.yaml.org) explicitly states that YAML stands for YAML Ain’t Markup Language and that “YAML is a human-friendly data serialization language for all programming languages.” So, YAML is primarily intended as a data serialization language, with the added goal of being as human-friendly as possible.

If you compare YAML to the other data formats like XML or JSON, it seems to do much the same thing: it represents constructs like lists, key-value pairs, strings, and integers. However, as you’ll soon see, YAML does this is a uniquely human-readable way. YAML is very easy to read and write once you understand how its syntax maps to these basic data structures.

This is a big reason that many automation tools use YAML as a method of defining an automation workflow or providing a data set to work with (like a list of VLANs). YAML also has the added benefit of helping to enable IaC approaches, covered in Chapter 13.

At the time of this writing, the latest YAML specification is YAML 1.2.2, published at https://www.yaml.org. Also provided on that site is a list of software projects that implement YAML, typically for the purpose of being read into language-specific data structures and doing something with them. If you have a favorite language, it might be helpful to follow along with the YAML examples in this chapter and try to implement them using one of these libraries.

Let’s take a look at some examples. Let’s say you want to use YAML to represent a list of network vendors. If you paid attention in the preceding section, you’re probably thinking that you want to use a string to represent each vendor name—​and you’d be correct! This example is simple:

- Cisco
- Juniper
- Brocade
- VMware
This YAML document contains four items. You know that each item is a string. One of the nice features of YAML is that you usually don’t need quotes or double quotes to indicate a string; a string is usually automatically discovered by the YAML parser (e.g., PyYAML). Each of these items has a hyphen in front of it. Since all four of these strings are shown at the same level (no indentation), you can say that these strings compose a list with a length of 4.

YAML closely mimics the flexibility of Python’s type system. A good example of this flexibility is shown by mixing data types in a list:

- Core Switch
- 7700
- false
- ['switchport', 'mode', 'access']
This is another list, again with a length of 4. However, each item is a totally unique type. The first item, Core Switch, is a string type. The second, 7700, is interpreted as an integer. The third is interpreted as a boolean. This interpretation is performed by a YAML interpreter, such as PyYAML. PyYAML, specifically, does a pretty good job of inferring the kind of data the user is trying to communicate.

The fourth item in this example is itself a list, containing three string items. This is an example of a nested data structure in YAML. You’ve also seen the various ways that some data can be represented. Our “outer” list is shown on separate lines, with each item prepended by a hyphen. The “inner” list is shown on one line, using brackets and commas. These are two ways of writing the same thing: a list.

TIP
Sometimes it’s possible to help the parser figure out the type of data you wish to communicate. For instance, if you want the second item to be recognized as a string instead of an integer, you can enclose it in quotes ("7700"). You also enclose data in quotes if a string contains a character that is part of the YAML syntax itself, such as a colon (:). Refer to the documentation for the specific YAML parser you’re using for more information.

Early on in this chapter, we briefly talked about key-value pairs (or dictionaries, as they’re called in Python). YAML supports this structure quite simply. Let’s see how you might represent a dictionary with four key-value pairs (Example 8-3).

Example 8-3. YAML dictionary with mixed types
Juniper: Also a plant
Cisco: 6500
Brocade: True
VMware:
  - esxi
  - vcenter
  - nsx
Here, your keys are shown as strings to the left of the colon, and the corresponding values for those keys are shown to the right. If you want to look up one of these values in a Python dictionary, for instance, you reference the corresponding key for the value you are looking for.

Similar to lists, dictionaries are flexible with respect to the data types stored as values. In Example 8-3, you are storing a myriad of data types as the values for each key-value pair.

YAML dictionaries—​like lists—​can be written in multiple ways. From a data representation standpoint, the previous example is identical to this:

{Juniper: Also a plant, Cisco: 6500, Brocade: true,
VMware: ['esxi', 'vcenter', 'nsx']}
Most parsers will interpret these two YAML documents precisely the same, but the first is obviously far more readable. The latter is a good illustration of the close relationship between YAML and JSON, but from a practical perspective, you’ll rarely need to use the latter format. Again, the primary use case for YAML is to be human readable, so stick with the conventions that most closely align with this.

Finally, you can use a hash sign (#) to indicate a comment. This can be on its own line or after existing data:

- Cisco    # ocsiC
- Juniper  # repinuJ
- Brocade  # edacorB
- VMware   # erawMV
Anything after the hash sign is ignored by the YAML parser.

As you can see, YAML offers a friendly way for human beings to provide structured data to software systems. However, YAML is fairly new as far as data formats go. When it comes to data formats used for communication directly between software elements (i.e., no human interaction), other formats like XML and JSON are much more popular and have much more mature tooling that is conducive to that use case.

Working with YAML in Python
Let’s narrow in on a single example to see exactly how a YAML interpreter will read in the data you’ve written in a YAML document. Let’s reuse one of the previous examples to illustrate the various ways to represent certain data types:

Juniper: Also a plant
Cisco: 6500
Brocade: true
VMware:
  - esxi
  - vcenter
  - nsx
Let’s say this YAML document is saved to your local filesystem as example.yml. Your objective is to use Python to read this YAML file, parse it, and represent the contained data as some kind of variable.

Fortunately, the combination of native Python syntax and the aforementioned third-party YAML parser, PyYAML, makes this easy:

import yaml
with open("example.yml") as f:
    result = yaml.load(f)
    print(result)
    type(result)

{'Brocade': True, 'Cisco': 6500, 'Juniper': 'Also a plant',
'VMware': ['esxi', 'vcenter', 'nsx']}
<type 'dict'>
TIP
The Python snippet in the preceding example uses the yaml module that is installed with the PyYAML Python package. This is easily installed using pip as discussed in Chapter 6.

This example shows how easy it is to load a YAML file into a Python dictionary. First, a context manager is used to open the file for reading (a common method for reading any kind of text file in Python), and the load() function in the yaml module allows us to load this directly into a dictionary called result. The lines that follow this code show that this has been done successfully.

XML
As mentioned in the previous section, while YAML is a suitable choice for human-to-machine interaction, other text-based formats like XML and JSON tend to be favored when software elements need to communicate with one another. This section covers Extensible Markup Language (XML), why it is suitable for this use case, and some of the ecosystem tools that exist for working with it.

The XML specification is defined and maintained by the World Wide Web Consortium, or W3C. XML was derived from a similar but older format called Standard Generalized Markup Language (SGML). XML is considered a subset of SGML, and as a result, any existing SGML parsers should be able to parse XML.

XML was originally created in the late 1990s, when the World Wide Web was moving from static HTML pages to more dynamic content that required lightweight update mechanisms. During this time, the limitations of HTML on its own in this respect were becoming obvious. HTML was designed for the sole purpose of describing the format and structure of a web page, and as a result was quite static and not very extensible. XML was created so that arbitrary data—not just web-focused markup—could easily be transmitted over the network. Some of the earliest use cases for XML were applied toward creating a more dynamic web, but XML itself is a generic format for representing just about anything.

An early popular use case for XML was in the implementation of Asynchronous JavaScript and XML, or Ajax. This was one of the first web development techniques for making web content more dynamic. It accomplished this by having web applications send and receive data in the background, and use this data to dynamically refresh components within the application, without requiring a full page refresh. Another popular use case was SOAP, which was an RPC technique based on XML. At the time of this writing, both use cases have been supplanted by more modern, lightweight alternatives.

In the world of modern network automation, the most popular use case for XML is within the NETCONF protocol. In addition, while JSON is generally a more popular option, XML can be used as the data format for REST APIs as well. We talk about both of these in Chapter 10.

XML shares some similarities with YAML. For instance, it is inherently hierarchical. We can easily embed data within a parent construct, as shown in Example 8-4.

Example 8-4. Basic XML document
<device>
  <vendor>Cisco</vendor>
  <model>Nexus 7700</model>
  <osver>NXOS 6.1</osver>
</device>
In this example, the <device> element is said to be the root. While spacing and indentation don’t matter for the validity of XML, you can easily see the root, as it is the first and outermost XML tag in the document. It is also the parent of the elements nested within it: <vendor>, <model>, and <osver>. These are referred to as the children of the <device> element, and they are considered siblings of one another. This structure is conducive to storing metadata about network devices, as you can see in this particular example. An XML document may contain multiple instances of the <device> tag (or multiple <device> elements), perhaps nested within a broader <devices> tag.

You’ll also notice that each child element contains data within. Whereas the root element contains XML children, these tags contain text data. Thinking back to the section on data types, it is likely these would be represented by string values in a Python program, for instance.

XML elements can also have attributes:

<device type="datacenter-switch" />
When a piece of information has associated metadata, it may not be appropriate to use a child element to describe that metadata, but rather an attribute. Of course, you can do both if needed. The key is to understand the difference between data and metadata (data about data) and use the appropriate tool to describe it.

An XML document can contain tags with just about any kind of name, depending on the use case. You could, therefore, encounter a naming conflict when creating tags for your own XML data structure. For instance, you might choose to use the tag <device> to describe one of those fancy new “smartphones”:

<device>Palm Pilot</device>
However, what if you also want to use the tag <device> to describe a ToR switch? Fortunately, the XML specification has implemented a namespace system, which helps disambiguate collisions like this. XML allows you to define these namespaces, and refer to them using the xmlns attribute (Example 8-5).

Example 8-5. XML namespaces
<root>
  <e:device xmlns:c="https://example.org/enduserdevices">Palm Pilot</e:device>
  <n:device xmlns:m="https://example.org/networkdevices">
    <n:vendor>Cisco</n:vendor>
    <n:model>Nexus 7700</n:model>
    <n:osver>NXOS 6.1</n:osver>XML Schema Definition
  </n:device>
</root>
The basic primitives of XML are quite simple. However, to do something meaningful with XML, you should look at the tools available for working with XML in a programming language like Python.

Working with XML in Python
Python includes native support for searching and creating XML documents in its standard library, under the xml module. Popular third-party libraries, such as lxml, offer a similar API but different underlying implementation. For the sake of simplicity, we stick with what’s available natively in Python for these examples.

XML is inherently hierarchical, which makes it a good fit for a tree structure. This is made a bit more apparent in Figure 8-3, which provides a visual representation of the basic XML document in Example 8-4.

npa2 0803
Figure 8-3. Visualization of an XML tree structure
Visualizations like these can be helpful when navigating XML documents by using a programming language like Python.

There are a few ways to search a tree for a particular piece of information, each with its own pros and cons. In this section, we stick with whatever’s most pragmatic and straightforward to understand for a given example.

NOTE
This section is an extremely condensed explanation of tree structures and how to iterate through them. Courses and other learning resources focused on data structures and algorithms will give you a deeper understanding of the various kinds of tree structures and what they’re used for.

The API for working with XML in Python is based on the concept of an element tree: literally, a tree of XML elements. This is in fairly stark contrast to the way you work with data from YAML or JSON, which maps directly into data structures like dictionaries and lists, which can be easier to understand at first glance. However, if you’re able to understand the very basics of tree structures, the ElementTree API should be fairly straightforward.

You can import the ElementTree class directly from the standard library. You’re also creating the simpler alias ET so you can easily refer to it in the following examples:

>>> import xml.etree.ElementTree as ET
There are a few ways to import an XML document, such as with ET.parse(), which loads from a file on the filesystem. However, you can also load XML from a string variable, which you might have if you’re looking at a response to an API request. In the next example, you’re declaring your own variable data and then using the ET.fromstring() method to read this string and create a new element tree from it:

data = """
<devices>
    <device name="sw01">
        <vendor>Cisco</vendor>
        <model>Nexus 7700</model>
        <osver>NXOS 6.1</osver>
    </device>
    <device name="sw02">
        <vendor>Arista</vendor>
        <model>Arista 7800</model>
        <osver>EOS 4.27</osver>
    </device>
    <device name="sw03">
        <vendor>Juniper</vendor>
        <model>QFX 10008</model>
        <osver>Junos 21.3</osver>
    </device>
</devices>
"""

tree = ET.fromstring(data)
You can simply print the value of your new tree variable to see what has been created for you:

>>> print(tree)
<Element 'devices' at 0x7f953cc8e1d0>
You’ll notice that the type for this variable is Element. This doesn’t fully represent the whole tree, but only the root element, which is the outermost tag in your XML document: <devices>. This element will have references to its children that you can access, and those will have their own children, and so on.

The main way to access the children for a given Element is through iteration. You can create a for loop to iterate over this element, and the items provided at each iteration will be one of that element’s children:

>>> for device in tree:
...     print(f"Device {device} found!")
...
Device <Element 'device' at 0x7f953cc3c590> found!
Device <Element 'device' at 0x7f953cbeb6d0> found!
Device <Element 'device' at 0x7f953cbeb860> found!
Of course, this doesn’t tell you much about each device, only that there are three of them. To access more information, you must go a bit deeper, since the elements like model and vendor are child elements of the device elements. You can use the find() method to search within the children of a given node and find the first one that matches a given tag:

>>> for device in tree:
...     model = device.find('model').text
...     print(f"Device model is {model}")
...
Device model is Nexus 7700
Device model is Arista 7800
Device model is QFX 10008
Since you were already iterating through the children of tree to get the device elements, the find() method can be used directly on a device element to search its children.

In some cases, especially when you have a deeply nested structure and are trying to get to a particular element, the iter() method can be useful. It allows you to iterate over all tree elements with a certain tag:

>>> for vendor in tree.iter('vendor'):
...     print(vendor.text)
...
Cisco
Arista
Juniper
This saves you from having to use things like nested for loops, or chained find() or findall() calls. In this case, you can just ask for all vendor elements in the entire tree and iterate over them.

As you might imagine, searching large trees can get a little complicated. Fortunately, Python includes limited support for XPath, an expression language that helps simplify searching through an XML document. You can provide a simple XPath expression as a parameter to the findall() method, which will then return all elements in the tree that match that expression:

>>> for model in tree.findall("./device/model"):
...     print(model.text)
Nexus 7700
Arista 7800
QFX 10008
You can also locate nodes based on a combination of their metadata attributes (e.g., the name attributes) as well as their element name. This is where things can get really powerful. Let’s say you want to look up the model for the device named sw01:

>>> tree.find("./device[@name='sw01']/model").text
'Nexus 7700'
You built on the previous expression by specifying the attribute and the desired value alongside device in the path. However, you are still able to get a handle on the model element, because of the remaining /model portion of the expression.

A lot more remains to dig into here, but a detailed explanation of XPath is a bit outside the scope of this section. That said, if you’re planning to work with XML frequently, especially large XML documents, XPath is an important tool to have in your toolbox. The official Python documentation on the xml module contains a lot of helpful examples.

TIP
Another tool for working with XML data is XQuery. XQuery is a full-blown query language, similar to what SQL is for relational databases, whereas XPath is a way of providing simple expressions—typically, one-liners for locating data within an XML document. Because of their similarities, you may wonder whether you should learn both. XQuery is rarely needed in the context of network automation. Most of the time, a little bit of Python and XPath can get you just about anything you might need.

While XML is well represented throughout the history of network automation, you should be aware of at least one other text-based format, as we’ll explore next.

JSON
JavaScript Object Notation (JSON) is the final text-based data format we’ll look at, and arguably the most widely used. XML has seniority, and YAML fills a need for a human-readable format, but when it comes to the format chosen for transporting structured data within networked applications, especially those that use HTTP, JSON is the undisputed champion. Many of the tools and libraries that communicate with network APIs, such as those in Chapters 10 and 12, use JSON to send/receive structured data over the network.

The origin story of JSON is similar to that of XML, albeit slightly more recent. It too was created as a lightweight mechanism for exchanging data over the web, to enable more dynamic content. It was based on a subset of the JavaScript programming language (thus the name), and the types used within the JSON specification closely reflect those within JavaScript. However, JSON is a language-independent format that’s well supported by a multitude of programming languages and automation tools. Languages like Go and Python have native JSON support in their standard libraries.

JSON has gone through a few iterations when it comes to standardization but remains largely consistent with the original ECMA-404 standard even to this day. The current version of the Internet Standard for JSON is described in RFC 8259. As you can see from this, JSON is a remarkably simple format; RFC 8259 is only 16 pages long!

JSON can also be used for configuration-related use cases. Node Package Manager (npm) uses JSON to describe the configuration of an npm package. Cloud providers like AWS and GCP use JSON files to configure a variety of their command-line utilities. JSON does have subjective advantages over YAML in this regard. Unlike YAML, JSON does not use indentation to indicate the scope of a given block of data, but rather the more explicit curly brace ({}) and square bracket ([]) syntax. This can make it easier to read and edit JSON documents for those unaccustomed to using indentation for scoping, as is done in programming languages like Python. However, this is almost entirely a matter of preference.

NOTE
JSON is widely considered a subset of YAML. In fact, many popular YAML parsers can also parse JSON data as if it were YAML (you may recall some of the “alternative” syntax we used in that section, which is remarkably similar to JSON). However, some of the details of this relationship are a bit more nuanced. See the YAML specification for more information.

When compared directly against XML, it’s easy to see that JSON is more lightweight; it is generally able to describe the same underlying data with less overall text structure. Let’s say you want to represent a list of book authors in XML. You might do it like this:

<authors>
    <author>
        <firstName>Christian</firstName>
        <lastName>Adell</lastName>
    </author>
    <author>
        <firstName>Scott</firstName>
        <lastName>Lowe</lastName>
    </author>
    <author>
        <firstName>Matt</firstName>
        <lastName>Oswalt</lastName>
    </author>
</authors>
Example 8-6 shows the equivalent data structure in JSON.

Example 8-6. Equivalent JSON
{
    "authors":[
        {
            "firstName": "Christian",
            "lastName": "Adell"
        },
        {
            "firstName": "Scott",
            "lastName": "Lowe"
        },
        {
            "firstName": "Matt",
            "lastName": "Oswalt"
        }
    ]
}
You can see that JSON is clearly a more lightweight way of representing data. This results in a more efficient way of transmitting the same underlying data. Especially in the early 2000s, this had a meaningful impact on web performance.

JSON has a fairly straightforward set of built-in types that are similar to those you might find in most programming languages, with some minor terminology differences. You’ll find they map to our YAML experience quite nicely:

Number
A signed decimal number.

String
A collection of characters, such as a word or a sentence.

Boolean
True or False.

Array
An ordered list of values enclosed in square brackets, []; items do not have to be the same type.

Object
An unordered collection of key-value pairs; keys must be strings (enclosed in curly braces, {}).

Null
Empty value. Uses the word null.

In Example 8-6, you can see several of these types in use. You’ll notice that the whole document is wrapped in curly braces. This means that the outermost (or root) type is an object, which contains key-value pairs.

NOTE
A JSON document that uses an object as its root, or outermost, type isn’t uncommon, but also isn’t the only option. The outermost type could also be an array, containing elements of any type.

In this example, the object being described contains only a single key-value pair (note that the keys within a JSON object are always strings). The key is authors, and the value for that key is an array. This is also equivalent to the list format we discussed in YAML—an ordered list of zero or more values. This is indicated by the square brackets [].

Contained within this list are three objects (separated by commas and a newline), each with two key-value pairs. The first pair describes the author’s first name (key of firstName) and the second, the author’s last name (key of lastName).

Working with JSON in Python
JSON enjoys wide support across a myriad of languages. A JSON document can often be mapped directly into native data structures in languages like Python (dictionaries, lists) and Go (slices, maps, structs). We’ll now look more specifically at how to work with JSON in Python.

Our JSON data is stored in a simple text file:

{
  "hostname": "CORESW01",
  "vendor": "Cisco",
  "isAlive": true,
  "uptime": 123456,
  "users": {
    "admin": 15,
    "storage": 10,
  },
  "vlans": [
    {
      "vlan_name": "VLAN30",
      "vlan_id": 30
    },
    {
      "vlan_name": "VLAN20",
      "vlan_id": 20
    }
  ]
}
Python has tools for working with JSON built right into its standard library, aptly called the json package. In Example 8-7, you can read this JSON file, convert it (load) into a Python dictionary, and print out some useful information about it (the inline comments can help explain each step in a bit more detail):

Example 8-7. Importing JSON to a Python dictionary
# Python contains very useful tools for working with JSON, and they're
# part of the standard library, meaning they're built into Python itself.
import json

# We can load our JSON file into a variable called "data"
with open("json-example.json") as f:
    data = f.read()

# Since our JSON document is an Object, json.loads() returns a dictionary.
# If our document was an Array, this would result in a list.
json_dict = json.loads(data)

# Printing information about the resulting Python data structure
print("The JSON document is loaded as type {0}\n".format(type(json_dict)))
print("Now printing each item in this document and the type it contains")
for k, v in json_dict.items():
    print(
        "-- The key {0} contains a {1} value.".format(str(k), str(type(v)))
    )
Those last few lines show exactly how Python views this data once it’s imported. The output that results from running this Python program is shown in Example 8-8.

Example 8-8. Results of importing JSON to a Python dictionary
~ $ python json-example.py

The JSON document is loaded as type <type 'dict'>

Now printing each item in this document and the type it contains
-- The key uptime contains a <type 'int'> value.
-- The key isAlive contains a <type 'bool'> value.
-- The key users contains a <type 'dict'> value.
-- The key hostname contains a <type 'unicode'> value.
-- The key vendor contains a <type 'unicode'> value.
-- The key vlans contains a <type 'list'> value.
NOTE
You might be seeing the unicode data type for the first time. In Python, the str type is just a sequence of bytes, whereas unicode specifies an actual encoding. The reason you’re seeing it here is that the JSON specification requires text to be encoded in Unicode. So, if you’re new to text encoding, you can conceptually think of unicode as a specific type of string and useful for the same kind of things as the str (string) type, discussed in Chapter 6.

Now that you’ve imported your JSON document into a native Python data structure, all the tools and techniques you learned in Chapter 6 can be used to find whatever information you’re looking for.

You can also perform the reverse action—that is, taking a Python data structure and creating a JSON document from it. Example 8-9 creates a Python list vendors and uses the method json.dumps() to create a JSON document containing an array.

Example 8-9. Dumping a Python list as a JSON array
>>> import json
>>> vendors = []
>>> vendors.append("Cisco")
>>> vendors.append("Arista")
>>> vendors.append("Juniper")
>>> print(json.dumps(vendors, indent=2))
[
  "Cisco",
  "Arista",
  "Juniper"
]
## Binary Data Formats
So far in this chapter, we’ve been discussing text-based data formats. These data formats leverage an intermediate representation for enhanced portability. The vast majority of languages and tools can easily understand formats like JSON. Therefore, when you’re sending data over the network, you can use this format and know that the other end can understand it, even if it’s written in another language, by another team, on another continent.

Hopefully, it has also been clear that while this intermediate step sacrifices a bit of efficiency to achieve this portability, the vast majority of use cases in network automation simply do not have the performance requirements for this inefficiency to become a problem. That said, in a few situations, the extra time and storage/bandwidth required to serialize and deserialize into these formats is considered prohibitively inefficient. For these, a more efficient type of format is called for: a binary data format. To understand this format, we first have to explore data types a bit more and how they actually work under the hood.

The type system in any statically typed programming language can usually be thought of as a set of aliases that the compiler uses to represent various lengths of bytes. For example, an int in Go, which defaults to a 32-bit integer, is just a way of allocating 4 bytes of memory. It follows, then, that if we were to define our own type composed of fields like these, that type would occupy as much memory as the sum of its constituent parts; see Example 8-10.

Example 8-10. A simple Go struct
type Coords struct {
    X int32
    Y int32
}
This type is another way of saying “a 64 bit chunk of memory.” The fields used within this type are a way of telling the compiler that the first 32 bits are used for one purpose and the second 32 bits for another.

As you might imagine, this mapping is language-specific. The way Go maps its own type system into memory is very different from the way Python, Rust, or C does it. However, applications written in these languages still need to communicate somehow. We’ve established that applications can serialize into one of the aforementioned text-based formats to solve this problem, but as we’ve also shown, this comes at a computational and storage cost that can be unacceptable for some use cases. This is where binary data formats come in.

The primary difference between binary data formats and the text-based formats you’ve seen thus far is that text-based formats require two separate steps to store or transmit information (encoding/decoding and serialization/deserialization), whereas binary data formats do it all in a single step. This is because applications using binary data formats do not have to first serialize data to an intermediate format like JSON or XML, but rather, they are able to read the raw bytes off of a network request or file on the filesystem, and map those bytes directly into a data structure, as shown in Figure 8-4.

npa2 0804
Figure 8-4. Decoding/deserialization and serialization/encoding binary data formats
NOTE
Of course, even JSON messages must eventually be translated into binary in order to be transmitted or stored; everything your computer does is eventually a bunch of 0s and 1s, after all. However, because binary data formats do not require a text-based serialization step, they require much less work to get to that stage, and ultimately require a smaller number of bits in which to do it. You can think of binary data formats as being “closer” to that low-level representation.

For a more networking-centric example of binary data formats, you need look no further than your earliest Cisco Certified Network Associate (CCNA) or CompTIA Network+ studies. Packet headers are laid out as raw bits and bytes. You are taught to understand that bits 64–72 of an IPv4 packet are for that packet’s time to live (TTL), and that the source and destination address are both 32 bits long. You know this because you spend time understanding the concepts that each field represents, and then you understand how many bits are required to represent them. Most importantly, this layout of bits is standardized and must be strictly adhered to; otherwise, two network endpoints wouldn’t be able to communicate effectively.

The binary data formats we’re talking about here are really no different: they’re intended to represent a variety of data structures in such a way that it can be transmitted/stored and encoded/decoded as efficiently as possible. Imagine for a second if this wasn’t the case—that the IPv4 packet header was some agreed-upon JSON object that first needed to be deserialized before it could be read by your router. What a preposterously inefficient design that would be! Why, then, is such an idea so patently absurd, yet it’s totally acceptable to use formats like JSON for our automation-related APIs? Why wouldn’t we just use binary data formats everywhere? Well, the reality is that some APIs have gone in this direction. As we discuss in Chapter 10, binary data formats are quite popular in newer technologies like gRPC. However, JSON-based APIs are hardly going away. And don’t forget, text-based data formats still have their unique advantages: they’re more readable, easier to debug, and enjoy broader application support.

Application developers have a slew of considerations when deciding between binary data formats and text-based data formats. However, given that this book is focused specifically on network automation, we can boil this down a bit. The first thing you should know is that not many use cases really require the efficiencies gained by using binary data formats. Most of the time when you’re running network automation workflows, the vast majority of your computer’s time is spent idle. Typically, your script or tool will send a few API requests and wait for a response. Using text-based data formats for use cases like this is perfectly fine.

Binary data formats are preferred by a network automation professional for two main reasons:

Developer preference
As you’ve seen, all data formats come with their own set of tools and techniques for working with them. Developers become accustomed to using certain tools, and barring a compelling reason to do otherwise, sticking with what you know is sometimes useful. In some cases, this manifests in the form of a vendor building an API that uses only a certain format; if you’re writing a script to work with this API, you too will have to understand how to work with this format.

Performance
This is a nuanced point that shouldn’t be taken for granted. Again, most of what is done in network automation is trivial for a computer to handle, However, some use cases, such as streaming network telemetry, do benefit from these kinds of efficiencies. The key questions to ask here are, “How much information is being transmitted, and how often is it being transmitted?” One of the reasons streaming telemetry often uses binary data formats is that it’s sending very frequent updates, each of which could contain a significant amount of data.

CAUTION
Binary formats aren’t always faster or more efficient than text-based formats. A lot of nuance can be found in the implementations behind each technology, and in how well the data itself is laid out in each format. It’s totally possible that a text-based format could outperform a binary format. While the scales are definitely tipped in favor of binary formats when it comes to performance, this outcome is not guaranteed by any means.

Because most of the specific benefits and drawbacks of binary data formats tend to vary based on the format in question, it’s time to dive into some examples. We have quite a few options when it comes to binary data formats, and we’ll start with Protocol Buffers.

Protocol Buffers
By far, the most likely binary data format you’ll run into in your network automation journey is that provided by Protocol Buffers, also known as protobuf. As described on its website, protobuf was originally developed by Google for internal use, as a smaller and faster alternative to XML.

Since then, the specification and tooling for protobuf have become open source, so that anyone can use them. Previously, if you wanted to store or transmit information between applications, you had to either pick one of the well-known text-based formats like XML (which can be very slow when you’re operating at Google scale), or use a binary format specific to a programming language like Python or Go, which solves the performance problem to some extent but then locks you into that format’s language and ecosystem.

Protobuf, in contrast, is language agnostic, with nine languages supported by the latest version of the specification, and more being added all the time. This is the result of two important and distinct components:

A schema definition language that allows you to specify services and messages in a way that’s not specific to any one application language (e.g., Python or Go)

Tools for automatically generating source code in any of the supported languages based on these definitions

We’ll explore both of these aspects in the following sections.

Protobuf definitions
Many text-based formats like JSON and XML are self-describing, in that they have formally defined textual conventions for indicating the type of data that they contain. Most binary data formats, on the other hand, are not self-describing. Because they’re just an opaque blob of bytes, you need some kind of external type definition so these bytes can be translated into types that both ends of a communication stream can make sense of.

Protobuf is no exception to this. To use protobuf to communicate, developers must first define their message format in the protobuf interface definition language, usually stored as .proto files. This language allows you to define the complex types that will eventually be serialized into binary data and transmitted to another application. Protobuf is a language-agnostic, human-editable way of representing otherwise opaque, binary data. It’s like a Rosetta Stone for binary data—and both sides of a protobuf-based communication channel must have the same copy of it in order to communicate.

A type in this language is referred to as a message. With this, you can define a custom type that you would like to be serialized via protobuf (Example 8-11).

Example 8-11. A protobuf message
message Router {
    int32 id = 1;
    string hostname = 2;
}
Before we go much further, a closer examination of the protobuf message in Example 8-11 is warranted:

message Router is a top-level declaration of the message named Router. You can use this message definition within other messages or as a parameter or return value for services, which we cover later.

There are two singular fields: id, which is of type int32, and hostname, which is of type string.

Each field has a number to the right that’s known as the field identifier. Since the order of serialized protobuf data is an implementation detail (and therefore can vary), this helps ensure that the raw binary is deserialized into the right fields.

At no point is the textual data in Example 8-11 transmitted on the wire (this is a key distinction between this and other text-based formats presented thus far, like XML or JSON). It is merely a human-editable specification that the protobuf tooling can read, and “compile” into a more efficient format. For the same reason you would define a class or a struct (as in Example 8-10) so your program knows how to allocate and read a chunk of memory on your computer, you can define a protobuf message in a .proto file so that your software knows how make sense of the raw binary data you’re retrieving or sending.

Protobuf messages can also reference other messages when the built-in types like int32 and string aren’t enough on their own. Let’s add a second message type called Interface, and then add a field to the Router message that uses it:

message Router {
    int32 id = 1;
    string hostname = 2;
    repeated Interface interfaces = 3;
}

message Interface {
    int32 id = 1;
    string description = 2;
}
The repeated keyword is the protobuf equivalent of a list or array. It indicates that the field interfaces is not just a single instance of the Interface message type, but multiple instances.

Protobuf also allows you to define a Service, which describes a set of RPC functions that can use either the built-in types, or the messages you’ve defined as parameters or return types:

service RouterService {
    rpc GetRouter(RouterRequest) returns (Router);
}
TIP
Service declarations can be used to define service endpoints for frameworks like gRPC and gNMI. For these, protobuf is one of the most popular data representation technologies, and Chapter 10 covers this in more detail.

This is an extremely light introduction to protobuf definitions. There are several more built-in types and important keywords to know if you want to write your own definitions or even read an existing definition. The Protocol Buffers Language Guide is a great next step if you want to dig in further.

Protobuf tooling and code generation
The second important component that makes protobuf work is the tooling that can automatically generate code in the language of your choice, from these message and service definitions. This automatically generated code makes it much easier to send and receive binary-encoded protobuf data that follows those definitions in your applications.

You looked at individual pieces of our protobuf definition in the preceding section, and a full working example is shown in Example 8-12.

Example 8-12. Full protobuf definition
syntax = "proto3";
package networkstuff;

service RouterService {
    rpc GetRouter(RouterRequest) returns (Router);
}

message RouterRequest {
    int32 id = 1;
}

message Router {
    int32 id = 1;
    string hostname = 2;
    repeated Interface interfaces = 3;
}

message Interface {
    int32 id = 1;
    string description = 2;
}
The number one tool you’ll want to familiarize yourself with, and make sure you have installed anywhere you want to write code that uses protobuf, is protoc. This is the protobuf compiler, and it allows you to go from the generic message definitions you’ve been working with thus far to “real” code that you can use in languages like Python or Go. Instructions for downloading and installing protoc can be found on the main Protocol Buffers website.

Once installed, protoc can be used to generate code for a variety of languages, even simultaneously. This can be executed on the bash command line. The following example instructs protoc to generate Go and Python code in the local directory, from the protobuf definition (as shown in Example 8-12) in the file networkstuff.proto:

protoc --go_out=. --python_out=. networkstuff.proto
This creates two files, one for Python (networkstuff_pb2.py), and one for Go (networkstuff.pb.go). These contain automatically generated type definitions and constructors for working with the language-specific implementation of the messages and services defined in our protobuf source file. You can then refer to these in your own code in order to use those types.

Let’s take a closer look at how to work with the generated Python code. You can open an interactive Python shell in the current directory and import the new module by name:

>>> import networkstuff_pb2
Within this module, each of your protobuf messages is given its own class, which you can instantiate:

>>> router = networkstuff_pb2.Router()
>>> router.id = 1337
>>> router.hostname = "r1"
It turns out the Python implementation is pretty smart. Normally, instances of Python classes permit the addition of attributes on the fly, but if you try to do that with your protobuf-generated class, you get an exception:

>>> router.foo = "bar"
Traceback (most recent call last):
  File "<stdin>", line 1, in <module>
AttributeError: Assignment not allowed (no field "foobar" in message object).
This adds a little bit of safety against things like “fat fingering” attribute names. It’s still Python, so this check doesn’t happen until runtime (with languages like Go, this would be caught at compile time), but it’s better than nothing.

You can also use the add() method on the router.interfaces attribute to instantiate a new Interface object:

>>> if1 = router.interfaces.add()
>>> if1.id = 1
>>> if1.description = "outside interface"
>>> if2 = router.interfaces.add()
>>> if2.id = 2
>>> if2.description = "inside interface"
Now that you have instantiated your protobuf-defined Router object in Python and populated it with sample data, you can use the SerializeToString() method to see the byte-level representation of this instance, printed as a Python byte string:

>>> router.SerializeToString()
b'\x08\xb9\n\x12\x02r1\x1a\x15\x08\x01\x12\x11outside interface\x1a\x14...'
You can write this binary data to the filesystem:

>>> f = open('serialized.bin', 'w+b')
>>> f.write(router.SerializeToString())
>>> f.close()
Back at the bash shell, you can then use a tool like hexdump to inspect the raw bytes in the file:

~$ hexdump serialized.bin
0000000 b908 120a 7202 1a31 0815 1201 6f11 7475
0000010 6973 6564 6920 746e 7265 6166 6563 141a
0000020 0208 1012 6e69 6973 6564 6920 746e 7265
0000030 6166 6563
0000034
Finally, to go full circle with this example, you can use protoc to decode these raw bytes back into a readable format. This decoding requires the original .proto file, as well as the name of the message you intend to decode and, of course, the binary file itself (passed via stdin):

~$ protoc --decode networkstuff.Router networkstuff.proto < serialized.bin

id: 1337
hostname: "r1"
interfaces {
  id: 1
  description: "outside interface"
}
interfaces {
  id: 2
  description: "inside interface"
}
The latest protobuf specification supports a canonical encoding in JSON, which can be really useful for working with systems that require a more traditional format. This way, you can primarily define your message types in protobuf, and serialize to binary data when you can, but still have the option to generate JSON from a given message when needed.

The method for producing this will vary based on the language, but the protobuf Python library contains a package for working with JSON from protobuf types:

>>> from google.protobuf.json_format import MessageToJson
>>> print(MessageToJson(router))
{
  "id": 1337,
  "hostname": "r1",
  "interfaces": [
    {
      "id": 1,
      "description": "outside interface"
    },
    {
      "id": 2,
      "description": "inside interface"
    }
  ]
}
Truthfully, as a network automation professional, you’re unlikely to use protobuf to write serialized binary data to the filesystem. Things get really exciting when you are able to build on what we’ve only started to explore here and leverage protobuf in modern RPC frameworks like gRPC. Technologies like this will also leverage the code generated by protoc not only for the message definitions, but also to generate functions that represent the services we defined. It doesn’t matter if you have an API client written in Python and an API server written in Go (or many other combinations); as long as both sides are working from the same protobuf definitions, they can communicate. Chapter 10 covers this in much more detail.

Protobuf is a modern, lightweight binary data format that in the world of network automation has already been widely adopted. However, before we move on, we should touch on a few other binary data formats that you may come across.

## Other Binary Data Formats
Within the scope of network automation, protobuf is really the only binary data format you need to know about, since it is a key component in many modern network programmability options. However, a few other binary data formats are potentially relevant. It’s useful to be aware of some of these alternatives, so we’ll spend a few sentences discussing a few of them and their pros and cons:

Pickle
Pickle is a binary format for serializing Python objects. It is specific to Python and therefore will not work in other languages, but offers support for serializing just about any kind of object structure you have in your Python programs. It has advantages like built-in de-duplication (will not serialize the same object twice) and backward compatibility.

Gob
Gob is a binary format for serializing Go types. It aims to have the same speed advantages as Protocol Buffers, but presented in a way that’s much easier to use and doesn’t require a separate interface definition language, as you would have in .proto files. You need only define your types in code, and the gobs package will be able to determine how best to serialize those types by using reflection techniques.

BSON
BSON is a “binary-encoded serialization of JSON-like documents.” It was originally invented as an internal representation of data for the MongoDB database. It’s more efficient than its textual counterpart JSON, but still somewhat less efficient than other binary formats, since it includes things like field names within the serialized data. It does include some additional types that are not supported in the JSON specification.

FlatBuffers
FlatBuffers is similar to protobuf, including the fact that both were originally developed at Google. However, unlike protobuf, FlatBuffers allows you to directly access the serialized data in the form of a flat binary buffer, without having to unpack or deserialize it first. You can also deserialize a portion of the buffer, as opposed to having to deserialize the entire buffer all at once. This is highly desirable for extremely performance-sensitive applications, such as video games.

Apache Thrift
Apache Thrift is also similar to protobuf, in that it is a binary data format that includes an RPC framework, an interface definition language, and code generation tooling. However, while it was originally created at Facebook, it has since become an Apache project. Thrift and protobuf are typically seen as roughly equivalent in terms of performance (most comparisons have these two tied at first place). Thrift does offer a full RPC implementation, whereas protobuf generates only RPC stub functions that need to be implemented to be useful.

It’s useful to be aware of these other formats, and each has its own benefits and drawbacks. However, as a network automation professional, the choice of which binary data format to use will almost always be made for you. A network platform will typically determine one of these and provide either message definitions for you to create your own code or a prebuilt library that you can simply consume.

Next, we’ll cover data modeling, which allows us to place additional constraints on the data sent using one of these formats.

Data Modeling
So far in this chapter, we’ve discussed a variety of data formats. Text-based formats like YAML, XML, and JSON are great for representing data in a human-readable and portable way. Binary data formats like protobuf are useful when performance is a bit more important. All these formats have basic type systems so your program is able to understand that a given series of characters or bytes is a string, integer, or boolean. At the end of the day, all these formats are aimed at representing data in a way that can be serialized and deserialized, to facilitate things like API-based communication.

However, sometimes we need more than just simple serialization. Let’s imagine that we are interacting with an API endpoint to update the hostname for a network device. The JSON payload for a request to this endpoint might look something like Example 8-13.

Example 8-13. Example JSON payload
{
  "hostname": ""
}
This JSON object has a single key, hostname, whose value is also a string—presumably representing the new hostname we want to use for this device. However, while it is a valid string, it is also empty. From a JSON formatting perspective, this is a perfectly valid syntax; any mature JSON parser will have no problem deserializing this document.

However, if we were to send this payload to the API endpoint in question, it could still cause problems. These problems would have nothing to do with the validity of the JSON document itself, but rather the downstream effects of sending an empty string as a parameter to the hostname update functionality that this API endpoint represents. Now, of course, the API server could include a check to ensure that this field is not empty; this might take the form of a conditional, as in Example 8-14.

Example 8-14. Explicit check for an empty string
req = json.loads(json_str)
if req["hostname"] == "":
    raise Exception("Hostname field must not be empty")
However, what about hostnames that are too long? What about special characters that might be supported in JSON but aren’t supported by the actual network device? By the way, all these considerations apply only to this one hostname field; what about even moderately more complex JSON payloads? We may have many more types in this payload to think about like integers, arrays, or nested objects, each with its own specific validity concerns. We might want to ensure that a given JSON array is not empty, or contains no more than five elements, or doesn’t contain any duplicates.

Writing server-side code to check for all these cases can quickly become unsustainable. Even if we could stay on top of all of them, such an approach would create an ugly experience for anyone writing code to consume such an API. If all these checks were built into the API server itself, clients would have great difficulty knowing how to send valid data to this API. The maintainers of the API would have to maintain detailed documentation about all of these checks and ensure that it was kept up-to-date manually (you can probably imagine how rarely this approach ends in success).

Data modeling is a set of tools and techniques for solving this problem. Whereas data formats allow you to serialize structured data generically, data modeling allows us to take this a step further and provide constraints that this structured data must adhere to. It gives us the opportunity to describe more specific rules and relationships that bring data into alignment with a specific use case or business process.

Typically, this is accomplished using some kind of data modeling language, which specializes in describing these constraints and relationships between data. These are often developed in conjunction with applications (i.e., an API server) designed to leverage that data model. Such an approach gives us key advantages and capabilities, especially within the network automation domain:

Data models can often be language agnostic. Multiple applications can use the same model (e.g., a Python API client and a Go API server). In addition, you have only one place to see or update the data model.

This approach focuses on the data, not the application. This makes it easier for nondevelopers (or developers who specialize in various languages) to understand the data model, without having to worry about language-specific syntax.

Many data-modeling techniques provide a way to generate application code that enforces these constraints. This can be extremely useful for working with APIs; if you have the data model, you can automatically generate code to reliably produce a correct payload for an API call.

The canonical example of data modeling in practice is the database schema, which is used to describe the organization and structure of data within a database. These schemas allow you to describe tables of data, which include columns of a particular data type, but also allow you to specify relationships between data, and constraints like the uniqueness of a particular value. For example, in relational database systems the primary key is a special designation describing a column of values that can be used to uniquely identify the row to which those values belong. Primary keys often enforce such a uniqueness constraint; an attempt to insert a new row with a primary-key value that already exists for another row will be rejected.

CAUTION
It may seem like we’re straying into software developer territory here. While it’s true that some of these concepts may be more aligned with the day-to-day work of professional developers, it doesn’t mean you’ll never need to create your own data models (even simple ones), and it certainly doesn’t mean that understanding the concepts behind data modeling, the technologies involved, or an existing data model aren’t profoundly useful skills to have as a network automation professional.

That said, more goes into creating a solid data model than throwing a few fields together. It can often require a more in-depth understanding of relationships between data, cardinality, and data normalization (or in some cases, denormalization), which can be tough for even seasoned software developers to get right.

As a result, this section doesn’t cover every aspect of data modeling, for every possible use case. Rather, we give you just enough insight into the important concepts involved with data models, some of the specific tools and techniques for creating and evolving them, and the most likely ways you’ll need to use this knowledge in your work as a network automation engineer.

While we won’t be diving into database schemas in this chapter, the idea of a schema as applied to the data formats we’ve discussed thus far is very much applicable. In its most general definition, a schema is just a way to describe the structure of data. As a result, we use the terms data model and schema somewhat interchangeably in this section, since they’re both close enough approximations for accomplishing our goals.

Before getting started, here are a few key points to keep in mind as you read the remaining sections of this chapter:

Data modeling involves the creation of a schema to which data must conform. This allows us to go beyond simple serialization and provide a more opinionated structure of the data that is relevant to our business logic or use case.

Data-modeling languages and tools are not serialization formats. They are not used to carry information, but only to describe it. You won’t see any of these modeling technologies in a packet capture or browser network trace.

Some data-modeling technologies are specific to a corresponding serialization format (e.g., JSON and XML), and others are a bit more broadly applicable.

We don’t cover every data modeling tool in existence—only those that you’re most likely to run into in your network automation journey.

Within the context of network automation, you can consider data models as analogous to a grammar textbook. It doesn’t tell you the specific words to say to a friend during a conversation, only the rules you should follow to ensure that the two of you can have a conversation of any kind. When you speak, you don’t regurgitate the textbook itself; rather, you use your own words that follow the rules from that textbook. In the same way, data models provide the specific rules and constraints that a particular communication mechanism must follow. Given that both sides are following the same “grammar textbook,” they can communicate.

YANG
Without a doubt, the data-modeling technology you’re most likely to run into during your network automation journey is YANG. Originally published as IETF RFC 6020, YANG was created as a data-modeling language specific to the NETCONF protocol, which we cover in greater detail in Chapter 10. However, in the most recent version, which is defined in RFC 7950, YANG has begun to decouple itself from NETCONF and XML so that other serialization formats like JSON, defined in RFC 7951, can be used, as well as other APIs like RESTCONF (NETCONF over an HTTP transport). Regardless, the main purpose of YANG is to model configuration and operational state data such as that transmitted during NETCONF RPCs.

You may have heard that a given NOS or API is YANG-based, or model driven using YANG. This is a way of summarizing an architectural approach to building programmable network systems that places the data model at the center. This is usually a good thing; starting with the data model allows a vendor to automatically generate code from that model to implement API servers, clients, and internal systems. It’s a much less fragile, less burdensome approach than building API bindings by hand.

CAUTION
One unfortunate by-product of condensing the YANG approach into such simple terms is that it almost sounds like YANG is used as a serialization format for APIs like RESTCONF or NETCONF. This is a popular misconception about YANG. YANG is not a serialization format like JSON or XML, and you won’t see YANG syntax in a packet capture of an API request or response. In fact, APIs that leverage YANG for data modeling usually use XML as the serialization format when sending data between a server and a client. For this reason, many of the examples in this section use XML to show how data that is modeled in YANG can be serialized in an API request.

YANG enjoys broad adoption by many companies and organizations. Many network vendors use YANG to build their systems with a model-driven approach. End-user-led organizations like OpenConfig aim to create a common set of vendor-neutral data models. The IETF also has working groups for building its own set of vendor-neutral models.

Like other data-modeling technologies, YANG enables you to define the constraints of that data—such as those found in a network configuration or state table. You can specify, for instance, that VLAN IDs must be between 1 and 4094. You can enforce the operational state of an interface, in that it must be “up” or “down.” Through these models, the behavior of data within and between network systems can be defined.

Various types of YANG models exist. Some of these YANG models were created by end users; others were created by vendors or open working groups:

Industry standard models include those from groups like the IETF and the OpenConfig Working Group. These models are vendor and platform neutral. Each model produced by an open standards group is meant to provide a base set of options for a given feature.

Of course, vendor-specific models also exist. Almost every vendor has its own solution for multichassis link aggregation groups (MC-LAGs), for example, each with its own variances in configuration and state data. As a result, each vendor would need to build a data model specific to these implementations.

Even within a single vendor, variances arise in the way a given feature is implemented across product platforms and would similarly require unique models.

As you may recall from earlier in this chapter, XML closely resembles a tree structure. Since YANG was originally intended to model data serialized in XML, it makes sense that the primitives it offers also follow this pattern. In fact, one of the core concepts in YANG is leaf, which allows you to define a singular piece of data that contains a single value and has no children. Note also the type statement, which allows you to specify that this element is a string, but other types are supported:

leaf hostname {
    type string;
    mandatory true;
    config true;
    description "Hostname for the network device";
}
This maps neatly to the XML document in Example 8-15.

Example 8-15. XML document satisfying the YANG model
<hostname>sw01</hostname>
This leaf statement is fairly flexible but still enforces some constraints on the data being described. For example, the mandatory true; statement means this field cannot be empty or blank. If you had omitted the sw01 text from Example 8-15, it would not validate against your YANG data model.

You may also remember that XML can contain multiple instances of the same element. A good example of this in practice is the list of configured DNS servers on a device. The leaf-list statement allows you to model this kind of data:

leaf-list name-server {
    type string;
    ordered-by user;
    description "List of DNS servers to query";
}
The ordered-by statement controls whether the order of elements within this data structure should be respected and maintained, or whether the implementation of the system can order the elements in the way it sees fit. The latter can be useful for things like VLAN definitions, as the order in which VLANs are defined doesn’t really matter. However, for other things like DNS name servers or access-list entries, order is extremely important. As a result, the statement ordered-by user; is used.

Again, here’s an example of XML data that adheres to this model:

<name-server>1.1.1.1</name-server>
<name-server>8.8.8.8</name-server>
Until now, we’ve been looking only at elements that don’t include any nested data. In YANG parlance, these are leaves in the tree. However, as we’ve shown, a nested structure is usually a more practical way of representing data like this. For instance, a specific VLAN may have several fields to describe it: minimally, a VLAN ID and a human-readable name. These could be represented as children of a generic vlan element:

<vlan>
    <id>100</id>
    <name>web_vlan></name>
</vlan>
<vlan>
    <id>200</id>
    <name>app_vlan></name>
</vlan>
YANG provides another way of defining lists, but unlike the leaf-list statement, the list statement is used when the elements of that list are themselves parent elements—that is, they contain nested elements, as shown in Example 8-16.

Example 8-16. List statement in YANG
list vlan {
    key "id";
    unique "name";
    leaf id {
        type int16;
    }
    leaf name {
        type string;
    }
}
This is also where you get to see some useful constraints in action. This key statement indicates that the id field should be used as a unique identifier for elements in this list. This is roughly analogous to a primary key in database terms. The unique statement specifies that the name value within these list elements should also be unique. This is useful for values that may not be used as a key, but should still be unique; in this case, it’s useful to ensure that the VLAN names are not duplicated.

However, you’re missing an important constraint here. The type for the VLAN ID is int16—which is quite broad when you consider that VLAN IDs are only positive and can go up to only a value of 4094 (16-bit signed integers can represent values from −32,768 to 32,767). Unfortunately, an 8-bit integer would be too small for this purpose. So, what do you do?

You can define your own custom data type that helps enforce these kinds of constraints. Using the typedef statement, you can define a new type by name—say, vlanid:

typedef vlanid {
    type int16 {
      range "1 .. 4094";
    }
}
Within this block, you can specify that this new type definition inherits from the built-in type int16, but then also enforces a constraint that the value must be within the range 1 to 4094. As a result, any element that references this type must not only be a 16-bit integer, but also fit within this more specific range.

You can then amend the leaf node from Example 8-16 to use this new type definition:

list vlan {
    key "id";
    unique "name";
    leaf id {
        type vlanid;
    }
    leaf name {
        type string;
    }
}
Of course, seeing a series of vlan elements at the root of an XML document would be strange. It’s more likely that these would be nested within a parent element like vlans:

<vlans>
    <vlan>
        <id>100</id>
        <name>web_vlan></name>
    </vlan>
    <vlan>
        <id>200</id>
        <name>app_vlan></name>
    </vlan>
</vlans>
This can be specified in YANG by using the container statement:

    container vlans {
        list vlan {
            key "id";
            unique "name";
            leaf id {
                type vlanid;
            }
            leaf name {
                type string;
            }
        }
    }
This was just a taste of some of the more common primitives within YANG syntax. YANG has plenty of other terms and concepts that we don’t have time to get into. The YANG RFCs are surprisingly readable and quite thorough, so these can be a reliable reference if you want to dive deeper. Numerous online resources are also available for YANG fundamentals and practical applications.

Before we end this section, however, it’s useful to take a look at ways you can do something practical with a YANG model by using some of the tools in its ecosystem. One of the most popular of these is pyang, which is a Python-based library as well as command-line tool for working with YANG. With pyang and some plug-ins built for it, you can do things like validate that a given model is compliant with YANG RFCs, check whether a given XML document is valid against a given model, and even generate a Python class hierarchy from a model.

Once pyang is installed, you can use it from the bash command line to check that a given YANG model is valid:

~$ pyang config.yang
If any validation errors result from parsing this YANG model, this command would list them in the resulting output.

pyang also allows you to convert a YANG model to several supported output formats. One useful output is tree, which provides a nice tree-like map of the module and its statements and types:

~$ pyang config.yang -f tree
module: config
  +--rw hostname        string
  +--rw vlans
  |  +--rw vlan* [id]
  |     +--rw id      vlanid
  |     +--rw name?   string
  +--rw name-servers
     +--rw name-server*   string
Another useful output format is a sample XML skeleton:

~$ pyang config.yang -f sample-xml-skeleton
<?xml version='1.0' encoding='UTF-8'?>
<data xmlns="urn:ietf:params:xml:ns:netconf:base:1.0">
  <hostname xmlns="https://example.org/config"/>
  <vlans xmlns="https://example.org/config">
    <vlan>
      <id/>
      <name/>
    </vlan>
  </vlans>
  <name-servers xmlns="https://example.org/config">
    <name-server>
      <!-- # entries: 0.. -->
    </name-server>
  </name-servers>
</data>
pyangbind is a plug-in for pyang that allows you to automatically generate a Python module from a YANG model. This can then be used to serialize out to XML, JSON, etc. You reference the plug-in directory where pyangbind can be located, and then refer to this location with the --plugindir flag. This makes a new output format available, pybind:

~$ PYANG_PLUGIN_DIR=$(pwd)/venv/lib/python3.8/site-packages/pyangbind/plugin
~$ pyang --plugindir $PYANG_PLUGIN_DIR -f pybind -o yangconfig.py config.yang
These commands generate a Python file, yangconfig.py, which you can then import from a Python prompt at the same location. Within this module is a class called config, which you can instantiate into cfg:

from yangconfig import config
cfg = config()
As of now, cfg is basically the Python equivalent of our YANG data model but is also empty. You can start populating it with information that is compliant with the model—for instance, a hostname:

cfg.hostname = "sw01"
Since the name servers are stored as a leaf list, these must be appended:

cfg.name_servers.name_server.append("1.1.1.1")
cfg.name_servers.name_server.append("8.8.8.8")
The VLANs are a bit more complicated, since these are modeled as a plain leaf. This means you need to add() a new VLAN by specifying its key as a parameter (this is the VLAN ID in our model). Then you can refer to it via that same key to set the other attributes, like name:

cfg.vlans.vlan.add(100)
cfg.vlans.vlan[100].name = "VLAN_100"
Note that you use a custom type to describe the VLAN ID, which specifies that it must be an integer between 1 and 4094. If you try to add a VLAN with an ID outside this range, an exception is raised:

>>> cfg.vlans.vlan.add(5000)

 (traceback omitted for brevity)

ValueError: 5000 does not match a restricted type
From here, you can serialize this into either XML or JSON. You can also use pyangbind to deserialize an existing XML or JSON document into this same class structure, and more. The pyangbind README file contains examples for these and other use cases.

As with most topics in this chapter, the preceding examples are really just a taste of what you can do with YANG. If you’re interested in providing a more structured, model-based approach to thinking about network data, you could do worse than to start off with YANG. However, you also should be aware of other modeling languages, which we’ll explore next.

### JSON Schema
JSON is an incredibly popular data format, especially in the frontend (web) developer world. As a result, it also enjoys a healthy ecosystem of tools and related specifications. JSON Schema is a data-modeling technology that allows you to easily document and validate JSON documents. If you know you want to use JSON as a data format, JSON Schema is a safe choice for creating a model or schema for validating the data you’re working with.

Like other data modeling technologies, JSON Schema includes a series of primitives and constraints for describing the layout of a set of data. As you might expect, its type system closely aligns with that of JSON itself; types like string, number, array, and object are all built right in. However, JSON Schema also provides a wide variety of other tools for describing the constraints within which those types should operate.

Let’s start with an example similar to the one you worked with for YANG, but instead of XML, you’ll use JSON; see Example 8-17.

Example 8-17. JSON document to model
{
    "hostname": "sw01",
    "vlans": [
        {
            "id": 100,
            "name": "VLAN_100"
        },
        {
            "id": 200,
            "name": "VLAN_200"
        }
    ],
    "nameservers": [
        "1.1.1.1",
        "8.8.8.8"
    ]
}
This JSON document contains an object type, which includes three fields:

hostname
This has a simple string value containing the device hostname.

nameservers
An array of strings containing our name servers.

vlans
An array of objects, which contain two fields, id and name. Each object represents a different VLAN.

As discussed in the previous section, you might want to make sure that the data shown in this example conforms to a few additional constraints beyond those imposed by the basic type system:

The vlans and nameservers arrays must not be empty, and they must not contain duplicate entries.

VLAN IDs must be between 1 and 4094.

All three fields—hostname, vlans, and nameservers—must be present; they cannot be omitted from the document.

A JSON Schema document is actually written in JSON, using a set of predefined terms and fields. Creating a new JSON Schema document starts by defining the outer type. Since our configuration data is a JSON object, you can specify this by using the type property. You can also provide useful metadata like title and description for our schema. The $schema property specifies the version of JSON Schema you’re using. This allows tools that use this schema to know which rules to use when parsing this schema and validating data with it. Example 8-18 shows a definition for a JSON schema.

Example 8-18. JSON Schema definition
{
    "$schema": "https://json-schema.org/draft/2020-12/schema",
    "title": "Config",
    "description": "A bit of configuration data for a network device",
    "type": "object",
    "properties": {
       ......omitted for brevity......
    },
    "required": ["hostname", "vlans", "nameservers"]
}
This definition includes two fields that you’ll see a lot more of in the following examples. The required key references keys within the JSON document that are mandatory; they should be present for a JSON document to be considered valid. The properties property allows you to specify further constraints that should be applied to these keys. Let’s now explore each property with specific examples (the following three examples are contained within properties).

The hostname key is fairly simple. You know it needs to be a string, but you also probably want to enforce a minimum and maximum length:

"hostname": {
    "type": "string",
    "minLength": 1,
    "maxLength": 20
}
The nameservers property is a bit more complex because it represents an array of values:

"nameservers": {
    "type": "array",
    "items": {
        "type": "string"
    },
    "minItems": 1,
    "uniqueItems": true
}
You need to not only specify the type array for this property but also use the items field to describe the type of the elements of that array—in this case, string. You can also specify that the array must have at least one item with the minItems field, and that the array shouldn’t contain duplicates by setting uniqueItems to true.

Finally, vlans is quite a bit more complicated since it is an array of objects:

"vlans": {
    "type": "array",
    "items": {
        "type": "object",
        "properties": {
            "id": {
                "type": "integer",
                "minimum": 1,
                "maximum": 4094
            },
            "name": {
                "type": "string"
            }
        },
        "required": ["id", "name"]
    },
    "minItems": 1,
    "uniqueItems": true
},
Fortunately, this part of the JSON document mostly uses terms you’ve already seen. The type of vlans is array, but its elements are of type object. Therefore, you also need to use the properties field to describe the properties of each object in the array. This is where you can use the required field again to specify that id and name are mandatory keys in each element. You can specify that id must not only be an integer, but also be between 1 and 4094.

Altogether, our JSON schema looks like this:

{
    "$schema": "https://json-schema.org/draft/2020-12/schema",
    "title": "Config",
    "description": "A bit of configuration data for a network device",
    "type": "object",
    "properties": {
        "hostname": {
            "type": "string",
            "minLength": 1,
            "maxLength": 20
        },
        "vlans": {
            "type": "array",
            "items": {
                "type": "object",
                "properties": {
                    "id": {
                        "type": "integer",
                        "minimum": 1,
                        "maximum": 4094
                    },
                    "name": {
                        "type": "string"
                    }
                },
                "required": ["id", "name"]
            },
            "minItems": 1,
            "uniqueItems": true
        },
        "nameservers": {
            "type": "array",
            "items": {
                "type": "string"
            },
            "minItems": 1,
            "uniqueItems": true
        }
    },
    "required": ["hostname", "vlans", "nameservers"]
}
You can store this schema as a JSON file, just as you can store the JSON document described in Example 8-18.

A plethora of tools exist for working with JSON Schema, in a variety of languages, for all kinds of use cases. One common use case is to simply validate that a JSON document adheres to a given schema. For instance, jsonschema is a popular Python-based tool for doing this. You can write Python scripts to use this library to perform validation, or you can use the command-line tool that comes with it to validate documents right on the bash command line:

~$ jsonschema --instance data.json schema.json
If there’s no output, you have a valid JSON document. However, you can easily tweak your JSON document to include invalid data to ensure that you’ve written a well-thought-out schema. For instance, say your document includes an invalid VLAN ID:

~$ jsonschema --instance data.json schema.json
50000: 50000 is greater than the maximum of 4094
Or say you’ve omitted the nameservers property:

~$ jsonschema --instance data.json schema.json
{
    'hostname': 'sw01',
    'vlans': [
        {'id': 100, 'name': 'VLAN_100'},
        {'id': 200, 'name': 'VLAN_200'}
    ]
}: 'nameservers' is a required property
As you’ve seen, JSON Schema can be a powerful tool for validating JSON data, and is probably a good choice if you know you’ll be working with JSON as a data format.

TIP
Since YAML and JSON are close relatives, some tools allow you to seamlessly validate YAML data as if it were JSON. In addition, it’s usually possible to convert a YAML document to JSON so that existing JSON-only tools can be used.

Next, we’ll explore how to validate data formatted in XML.

### XML Schema Definition
XML also has its own dedicated modeling language, known as XML Schema Definition, or XSD. One popular use case for XSD (as with most modeling languages) is to generate source code data structures that match and enforce the schema described by that data model. You can then use that source code to automatically generate XML that is compliant with that schema, as opposed to writing out the XML by hand.

For a concrete example of how this is done in Python, let’s look once more at our XML example:

<device>
  <vendor>Cisco</vendor>
  <model>Nexus 7700</model>
  <osver>NXOS 6.1</osver>
</device>
Your goal is to print this XML to the console by using some automatically generated code. You can do this by first creating an XSD document, and then using a third-party tool to generate Python code from that document. Then, that code can be used to print the XML you need.

Let’s write an XSD schema file that describes the data you intend to write out:

<?xml version="1.0" encoding="utf-8"?>
<xs:schema elementFormDefault="qualified" xmlns:xs="http://www.w3.org/2001/
XMLSchema">
  <xs:element name="device">
  <xs:complexType>
    <xs:sequence>
      <xs:element name="vendor" type="xs:string"/>
      <xs:element name="model" type="xs:string"/>
      <xs:element name="osver" type="xs:string"/>
    </xs:sequence>
  </xs:complexType>
</xs:element>
</xs:schema>
In this schema document, you are describing that each <device> element can have three children and that the data in each child element must be a string. Not shown here but supported in the XSD specification is the ability to specify that child elements are required; in other words, you could specify that a <device> element must have a <vendor> child element present.

You can use a tool called PyXB at the bash command line to create a Python file that contains class object representations of this schema:

~$ pyxbgen -u schema.xsd -m schema
This creates schema.py in this directory. So, if you open a Python prompt at this point, you can import this schema file and work with it. In Example 8-19, you’re creating an instance of the generated object, setting some properties on it, and then serializing it into XML by using the toxml() function.

Example 8-19. Generating XML from an XSD schema in Python
>>> import schema
>>> dev = schema.device()
>>> dev.vendor = "Cisco"
>>> dev.model = "Nexus"
>>> dev.osver = "6.1"
>>> dev.toxml("utf-8")
'<?xml version="1.0" encoding="utf-8"?><device><vendor>Cisco</vendor><model>Nexus
</model><osver>6.1</osver></device>'
Next, we’ll explore how to validate data that uses protobuf.

Modeling and Validating Protocol Buffers
Protocol Buffers don’t include any built-in data modeling or validation beyond their basic type system. If you want to do this, you need to look at third-party options. One popular choice is protoc-gen-validate, which is a plug-in to the protobuf compiler protoc that is maintained as part of the Envoy Proxy project.

This plug-in allows you to specify validation rules within the protobuf definitions we explored earlier. Let’s say you have a message Vlan with fields id and name:

message Vlan {
    int32 id = 1;
    string name = 2;
}
You can provide validation rules for this plug-in within brackets after the field number:

message Vlan {
    int32 id = 1 [(validate.rules).int32 = { gte: 1,  lte: 4094 }];
    string name = 2 [(validate.rules).message.required = true];
}
When compiled with this plug-in, the generated code will include validation methods on these types, such as Vlan.Validate(). This method can be used to determine whether a given class instance or struct adheres to the constraints described in these validation rules.

Summary
Data formats and data models are at the core of everything we do in network automation. Whether we’re talking about configuration management, troubleshooting, or even just generating quick reports, a firm grasp on these fundamentals is essential to being successful on your automation journey. As you’ve seen, specific technologies come and go, but the need for structured data, and the ability to describe the layout of that data and the constraints to which it must adhere, will never change.

Here are some parting thoughts:

Structured data is essential to a successful automation initiative. Unstructured data formats, while often ideal for humans to consume, are not designed to be easily parsed or understood by our automation systems and scripts.

None of the languages or formats discussed in this chapter are perfect. They’re all designed with specific trade-offs in mind; your job is to identify which trade-offs align best with your situation.

New data-modeling methods and languages are emerging all the time. For instance, CUE has recently grown in popularity as a bit of a hybrid between a schema definition language and a templating system. Some technologies discussed in this chapter have also decreased in popularity over time. This is a natural and expected evolution; keep your head on a swivel and assess each new tool on its own merits and trade-offs and how they align with your goals.

In the next chapter, we’ll use data from formats like those we’ve discussed here to drive the automatic creation of consistent, templated configurations.

# Templates
Much of a network engineer’s job involves the CLI, and much of this work requires syntax-specific keywords and phrases that are often repeated several times, depending on the change. This not only becomes inefficient over time but also is error prone. The way to configure a BGP neighbor relationship on Cisco IOS may be obvious, for instance, but what’s not obvious at times are the smaller, “gotcha” configurations, like remembering to append the right BGP community configuration. Often in networking, there are many different ways to do the same thing—and this may be totally dependent on your organization.

One of the key benefits of network automation is consistency—being able to predictably and repeatably make changes to production network infrastructure and achieve a desired result. One of the best ways to accomplish this is by creating templates for all automated interaction with the network.

Creating templates for your network configurations means that you can standardize those configurations for your organization, while also allowing network administrators and consumers (help desk, network operations center, IT engineers) to dynamically fill in some values when needed. You get the benefits of speed, requiring much less information to make a change, but also consistency (and through this, safety) because the template contains all the necessary configuration commands that your policies dictate.

This chapter starts with an introduction to template tools in general, and then presents specific implementations and shows how to leverage these tools to create network configuration templates.

The Rise of Modern Template Languages
Template technologies have been around for a very, very long time. Just a basic web search for “template languages” shows a multitude of these, most often several options for every related programming language.

You may also notice that the majority of these languages have deep applications in the web development industry. This is because much of the web is based on templates! Instead of writing HTML files for every single user-profile page that a social media site may have, the developers will write one and insert dynamic values into that template, depending on the data being presented by the backend.

In short, template languages have a wide variety of relevant use cases. Their obvious roots are in web development, and of course we’ll be talking about using them for network configuration in this chapter, but they have applications in just about any text-based medium, including documentation and reports.

So it’s important to remember that using templates requires three steps. First, the templates have to be written. Second, you need some form of data, which will ultimately get rendered into the template to produce something meaningful like a network configuration. This leads us to the third step: something has to drive data into the template. This could be an automation tool like Ansible, which we cover in Chapter 12, or you could be doing it yourself with a language like Python, which we show later in this chapter. Templates are not very useful on their own.

TIP
Most template languages aren’t full-on programming languages in the purest sense. Most often, a template language is closely tied to another language that will drive data into the templates that you’ve built. As a result, each template language and its “parent” language have several similarities. A good example is one that we heavily discuss in this chapter: Jinja is a template language that came out of a Python-centric community, so Jinja has distinct similarities with Python. So if you’re wondering which template language to use, it’s probably best to decide which “real” language you’re aligned with (either through writing your own code or by using an existing tool like Ansible) and go from there.

As mentioned previously, template languages aren’t necessarily a new concept, but we are seeing new ideas and even entire languages make it into the ecosystem all the time. If you look at the history of template languages, many were created to serve as a crucial part of the web: dynamic content. This is easily taken for granted these days, but back when the web was just getting started and most websites were built from fairly static content, dynamically loading pieces of data into a page was a big step forward.

Using Templates for Web Development
Django, a Python-based web framework, significantly leverages the concept of templated documents. Django has a template language that allows the web developer to create web content in much the same way they normally would, but also offers a way to make portions of the page dynamic. Using Django’s template language, the developer can designate portions of an otherwise static page to load dynamic data when the user requests a page.

Here’s a simple example—note that this looks much like an HTML document, but with certain portions replaced with variables (indicated with {{ }} notation):

<h1>{{ title }}</h1>

{% for article in article_list %}
<h2>
  <a href="{{ article.get_absolute_url }}">
    {{ article.headline }}
  </a>
</h2>
{% endfor %}
This template can be rendered by Django when a user loads the page. The Django framework will populate the title and article_list variables, and the user will receive a page that’s been fully populated with real data. The developer doesn’t have to write a static HTML page for every possible thing the user wants to retrieve; this is managed by logic on the backend of this web application.

NOTE
The Django templating language is similar (but not identical) to the templating language Jinja, which we discuss in depth in this chapter. Don’t worry about the syntax; we’ll get into that. For now, just focus on the concepts and the value that templates provide: consistency.

Detailing the multitude of other available template languages is outside the scope of this chapter, but you should be aware that they exist. Python alone has several options, such as the aforementioned Django and Jinja languages, but also Mako and Genshi. Other languages like Go and Ruby have built-in template systems. Again, the point to remember is that the important work of populating a template with data is the role of one of these languages, like Python or Go, so this is the number-one factor in deciding which template language to use. More often than not, it’s best to go with a template system built for that language.

Expanding On the Use of Templates
The concepts of templating, especially those discussed in this chapter, are not specific to any single use case and can be applied to nearly any text-based medium. At the end of the day, templates are just a way to transform data into a specific text format. Figure 9-1 illustrates this flow.

npa2 0901
Figure 9-1. How templates are produced
All template technologies discussed in this chapter generally work this way. This makes templates useful for anything text based, including reports, configuration files, and configurations. Perhaps you’re pulling data from a network device and want to be able to produce a nice report on this data and email it to coworkers. Example 9-1 shows a Jinja template for producing a report containing a list of VLANs.

NOTE
Full versions of the code examples in this chapter can be found in the book’s GitHub repo at https://github.com/oreilly-npa-book/examples/tree/v2/ch09-templates.

Example 9-1. Basic report with Jinja
| VLAN ID | NAME | STATUS |
| ------- |------| -------|
{% for vlan in vlans %}
| {{ vlan.get('vlan_id') }} | {{ vlan.get('name') }} | {{ vlan.get('status') }} |
{% endfor %}
Because you’re really just working with text, you can build a template for it. Keep the generic nature of templates in mind as you get into the details of specific use cases and particular template technologies like Jinja; templates have applications well beyond the narrow set of use cases presented in this chapter.

The Value of Templates in Network Automation
At this point, you might be wondering why we’re talking about web development and how that could possibly help you on your network automation journey. It’s important to understand the value behind templates in general, whether they’re used for the web or not. Templates provide consistency—instead of handcrafting text files full of HTML tags or entering CLI commands, templates enable you to declare which parts of your files need to remain static and which parts should be dynamic.

Every network engineer who has worked on a network long enough has had to prepare a configuration file for a new piece of network gear, like a switch or router. Maybe this has to be done for many switches—perhaps for a new data center build. To make the best use of time, it’s useful to build these configurations ahead of time, so that when the switch is physically racked and cabled, the network engineer needs only to paste the configuration into a terminal.

Let’s say you’re in charge of a rollout like this—it’s your job to come up with configurations for all switches going into the new data center being built for your organization. Obviously, each switch will need its own unique configuration file, but a large portion of the configuration will be similar among devices. For instance, you might have the same SNMP community strings, the same admin password, and the same VLAN configuration (at least for similar device types like ToR switches).

Then again, some parts of the configuration are probably unique to a single device. Routers or Layer 3 switches typically require unique addresses, and probably some fairly unique routing protocol configurations, depending on where that device exists in the topology. Deciding what parameters go to which switches can be time-consuming and likely to result in errors. Templates allow you to standardize on a common base configuration and help ensure that the right values get filled in for each device. Separating the template from the data that populates it is one of the best ways to simplify this process.

The primary value of templates for network engineers is achieving configuration consistency. Appropriately implemented, templates can reduce the likelihood of human errors caused when changing network configurations. There seems to be a lot of fear that making complex changes in an automated way in production is a bad idea, but if you follow good discipline and properly test your templates, you really can improve network operations. Templates don’t automatically remove human error, but when used properly, they can greatly reduce it, resulting in fewer outages.

Using templates to aid the rollout of new network devices is a simple yet powerful way to demonstrate the value of templates, since it has the added benefit of saving a lot of time for the network engineer. However, don’t think that this is the only place where templates can be used. In many network automation projects, templates are not even used by humans but by automation software like Ansible to push configuration changes to network devices—live, in production.

The following sections cover some of the most popular options for building template-based configurations and documents that exist for network automation.

Jinja
We mentioned Jinja in the introduction to this chapter, but it’s time to dive in a little deeper. While we’ll briefly explore a few other options in this chapter, you’ll notice that we spend a lot more pages on Jinja. Currently, Jinja is still by far the most popular template language in the world of network automation because of broad support by many automation tools and languages, including Python.

Let’s start with a basic example and write a template to configure a single switch interface. Here’s an example switchport configuration that you want to convert to a template (so you can configure the hundreds of other switchports in your environment):

interface GigabitEthernet0/1
 description Server Port
 switchport access vlan 10
 switchport mode access
This kind of snippet is easy to write a template for—you need only decide which parts of this configuration need to stay the same and which need to be dynamic. Example 9-2 has removed the specific interface name and converted it into a variable that you’ll populate when you render the template into an actual configuration.

Example 9-2. Parameterized interface name in a switchport config template
interface {{ interface_name }}
 description Server Port
 switchport access vlan 10
 switchport mode access
You can pass in the variable interface_name when rendering this template, and that spot will get filled in with the value associated with interface_name.

However, this template assumes that each network interface has an identical configuration. What if you want a different VLAN or a different interface description on some of the interfaces? In that case, you should also convert some of the other parts of the configuration into their own variables, as shown in Example 9-3.

Example 9-3. Several variables in the switchport config template
interface {{ interface_name }}
 description {{ interface_description }}
 switchport access vlan {{ interface_vlan }}
 switchport mode access
These are simple examples, but they’re not very namespace-friendly.

It’s common to leverage concepts like classes and dictionaries in a language like Python when rendering a template. This allows us to store multiple instances of data that we can loop over and write multiple times in our resulting configuration. We’ll look at loops in a future section, but for now, Example 9-4 shows that same template rewritten to take advantage of something like a Python class or dictionary.

Example 9-4. Using a dictionary in a switchport config template
interface {{ interface.name }}
 description {{ interface.description }}
 switchport access vlan {{ interface.vlan }}
 switchport mode access
This is a minor change, but an important one. The interface object is passed to the template as a whole. If interface is a Python class, then name, description, and vlan are all properties of that class. The same is true if interface is a dictionary—the only difference is that they are all keys of this dictionary, and not properties, so the rendering engine would automatically place the corresponding values for those keys when rendering this template.

Rendering a Jinja Template File in Python
In the previous example, you created a basic Jinja template for a switchport configuration, but you didn’t explore how that template is rendered or what drives data into the template to result in the final product. We’ll explore that now by using Python and the Jinja2 library.

NOTE
While the templating language itself is known as Jinja, the Python library for working with Jinja is called Jinja2.

Let’s use the same template snippet from Example 9-4 and use Python to populate those fields with real data. We’ll use the Python interpreter this time so you can walk through the example on your own machine.

TIP
The Jinja2 rendering engine in Python is not part of the standard library, so it is not installed by default. However, Jinja2 can be installed with pip3, through the command pip3 install jinja2, the same as any other Python package found on PyPI, as we covered in Chapter 6.

Once the Jinja2 library is installed, you should first import the required objects that you’ll need in order to render your templates:

>>> from jinja2 import Environment, FileSystemLoader
Next, you need to set up the environment, so the renderer knows where to find the template:

>>> ENV = Environment(loader=FileSystemLoader('.'))
>>> template = ENV.get_template("template.jinja")
The first line sets up the Environment object, specifying a single dot (.) to indicate that the templates exist in the same directory in which you started the Python interpreter. The second line derives a template object from that environment by statically specifying the template name, template.jinja. The contents of this file are identical to Example 9-4.

Now that this is done, you need your data, which you’ll specify using a dictionary. The keys for this dictionary correspond to the field names referenced in your template:

>>> interface_dict = {
...     "name": "GigabitEthernet0/1",
...     "description": "Server Port",
...     "vlan": 10,
...     "uplink": False
... }
CAUTION
You probably won’t have to manually create data structures in Python to populate a template with data. This is being done for illustrative purposes in this book. Typically, you’ll want to pull data from an API or a file on disk, rather than declaring it statically in your program or script.

You now have everything you need to render your template. You’ll call the render() function of your template object to pass data into the template engine, and use the print() function to output your rendered output to the screen:

>>> print(template.render(interface=interface_dict))
interface GigabitEthernet0/1
 description Server Port
 switchport access vlan 10
 switchport mode access
Note that you pass an argument to the render() function of your template object. Pay close attention to the name—the keyword argument interface corresponds to the references to interface within your Jinja template. This is how you get your interface dictionary into the template engine—when the template engine sees references to interface or its keys, it will use the dictionary passed here to satisfy that reference.

As you can see, the rendered output is as you expect. However, you don’t have to use a Python dictionary. It’s common to drive data from other Python libraries into a Jinja template, and this may take the form of a Python class.

Example 9-5 shows a Python program that’s similar to the code you just went through, but instead of using a dictionary, you use a Python class.

Example 9-5. Using a Python class to populate a Jinja template
from jinja2 import Environment, FileSystemLoader
ENV = Environment(loader=FileSystemLoader('.'))
template = ENV.get_template("template.jinja")

class NetworkInterface(object):
    def __init__(self, name, description, vlan, uplink=False):
        self.name = name
        self.description = description
        self.vlan = vlan
        self.uplink = uplink

interface_obj = NetworkInterface("GigabitEthernet0/1", "Server Port", 10)
print(template.render(interface=interface_obj))
The output from this program is identical to the previous output. Therefore, there really isn’t one right way to populate a Jinja template with data—it depends on where that data comes from. Fortunately, the Python Jinja2 library allows for some flexibility here.

Using Conditionals and Loops
It’s time to really make your templates work for you. The previous examples are useful for understanding how to insert dynamic data into a text file, but that’s just part of the battle of scaling network templates up to properly automate network configuration.

CAUTION
Jinja allows you to embed Python-esque logic into your template files in order to do things like make decisions or condense duplicate data into one chunk that is unpacked at render time via a for loop. While these tools are powerful, they can also be a slippery slope. Don’t get too carried away with putting all kinds of advanced logic into your templates. Jinja has some really useful features, but it was never meant to be a full-blown programming language, so it’s best to keep a healthy balance. Read the Jinja FAQ, specifically the section titled “Isn’t It a Bad Idea to Put Logic in Templates?” for some tips.

Using conditional logic to create a switchport configuration
Let’s continue our example of configuring a single switchport—but in this case, you want to decide what to render by using a conditional in the template file itself.

Often some switchport interfaces will be VLAN trunks, and others will be in mode access. A good example is an access layer switch, where two or more interfaces are the uplink ports that need to be configured to permit all VLANs. Our previous examples showed an uplink boolean property, set to True if the interface is an uplink and False if it is just an access port. You can check against this value in your template by using a conditional:

interface {{ interface.name }}
 description {{ interface.description }}
{% if interface.uplink %}
 switchport mode trunk
{% else %}
 switchport access vlan {{ interface.vlan }}
 switchport mode access
{% endif %}
TIP
If you run some of the examples in this section, you may notice extra whitespace, especially extra newlines between lines of rendered text. This is due to the way Jinja removes the template syntax when rendering the final text. You can use opt-in methods to manage this whitespace, such as placing a hyphen immediately following or preceding the substitution syntax (e.g., {% for item in seq -%}). This gives you the option to trim extra whitespace either before or after a given bit of Jinja logic. We’ve omitted these options to keep the examples simpler, but in practice you can and probably should use them to keep your rendered text cleaner, especially when rendering to formats where whitespace has a significant impact on the data being represented.

In short, if the uplink property of interface is True, you want to make this interface a VLAN trunk. Otherwise, let’s make sure that it’s set up with the appropriate access mode.

In the previous example, you also see a new syntax: the {% ... %} braces are a special Jinja tag that indicates some kind of logic. This template is built to configure GigabitEthernet0/1 as a VLAN trunk, and any other interface will be placed in access mode, in VLAN 10.

Using a loop to create many switchport configurations
You’ve configured only a single interface at this point, so let’s see if you can use Jinja loops to create configurations for many switchports. For this, you use a for loop that’s extremely similar to the syntax you would normally have in Python:

{% for n in range(10) %}
interface GigabitEthernet0/{{ n+1 }}
 description {{ interface.description }}
 switchport access vlan {{ interface.vlan }}
 switchport mode access
{% endfor %}
Note that you’re again using the {% ... %} syntax to contain all logic statements. In this template, you’re calling the range() function to give a list of integers to iterate over, and for each iteration, you print the result of n+1 because range() starts at 0, and normally, switchports start at 1.

Using a loop and conditionals to create switchport configurations
This creates an identical configuration for 10 switchports—but what if you want a different configuration for some of them? Take the example about Jinja conditionals—perhaps the first port is a VLAN trunk. You can combine what you’ve learned about conditionals and loops to accomplish this:

{% for n in range(10) %}

interface GigabitEthernet0/{{ n+1 }}
 description {{ interface.description }}
{% if n+1 == 1 %}
 switchport mode trunk
{% else %}
 switchport access vlan {{ interface.vlan }}
 switchport mode access
{% endif %}

{% endfor %}
This results in GigabitEthernet0/1 being configured as a VLAN trunk, but Gigabi⁠t​Ethernet0/2–10 is in access mode. Here is an example using simulated data for the interface descriptions:

interface GigabitEthernet0/1
 description TRUNK INTERFACE
 switchport mode trunk
interface GigabitEthernet0/2
 description ACCESS INTERFACE
 switchport mode access
interface GigabitEthernet0/3
 description ACCESS INTERFACE
 switchport mode access
...
Iterating over a dictionary to generate configurations
You were able to access keys in a dictionary in your Jinja template in the preceding section, but what if you want to iterate over dictionaries and lists by using a for loop? Let’s imagine that you’re passing the following list to your template as interface_list. Here’s the relevant Python:

intlist = [
    "GigabitEthernet0/1",
    "GigabitEthernet0/2",
    "GigabitEthernet0/3"
]
print(template.render(interface_list=intlist))
You then reference interface_list in your loop so that you can access its members and generate a switchport configuration for each one. Note that the nested conditional has also been modified, since your counter variable n no longer exists:

{% for iface in interface_list %}

interface {{ iface }}
{% if iface == "GigabitEthernet0/1" %}
 switchport mode trunk
{% else %}
 switchport access vlan 10
 switchport mode access
{% endif %}

{% endfor %}
You now simply refer to iface to retrieve the current item of that list for every iteration of the loop.

You can also do the same thing with dictionaries. Again, here’s a relevant Python snippet for constructing and passing in a dictionary to this Jinja template. Let’s keep it simple this time, and just pass a set of interface names as keys, with the corresponding port descriptions as values:

intdict = {
    "GigabitEthernet0/1": "Server port number one",
    "GigabitEthernet0/2": "Server port number two",
    "GigabitEthernet0/3": "Server port number three"
}
print(template.render(interface_dict=intdict))
You can modify your loop to iterate over this dictionary in much the same way you’d do in native Python:

{% for name, desc in interface_dict.items() %}
interface {{ name }}
  description {{ desc }}
{% endfor %}
The for name, desc... means that at each iteration of the loop, name will be a key in your dictionary, and desc will be the corresponding value for that key. Don’t forget to add the .items() notation as shown here, to properly unpack these values.

This allows you to simply refer to name and desc in the body of the template, and the result is shown here:

interface GigabitEthernet0/1
  description Server port number one

interface GigabitEthernet0/2
  description Server port number two

interface GigabitEthernet0/3
  description Server port number three
NOTE
In previous versions of Python, the order in which these interfaces appear would not be guaranteed, as dictionaries were by nature an unordered data structure. However, as of Python 3.7, dictionaries preserve insertion order for literals like these.

Generating interface configurations from a list of dictionaries
In this section, you’re going to combine using lists and dictionaries to really put this template to work for you. Each interface will have its own dictionary, and the keys will be attributes of each network interface, like name, description, or uplink. Each dictionary will be stored inside a list, which is what your template will iterate over to produce configuration.

First, here’s the data structure in Python that we just described:

interfaces = [
    {
        "name": "GigabitEthernet0/1",
        "desc": "uplink port",
        "uplink": True
    },
    {
        "name": "GigabitEthernet0/2",
        "desc": "Server port number one",
        "vlan": 10
    },
    {
        "name": "GigabitEthernet0/3",
        "desc": "Server port number two",
        "vlan": 10
    }
]
print(template.render(interface_list=interfaces))
This allows you to write a powerful template that iterates over this list, and for each list item, simply refers to keys found within that particular dictionary. Example 9-6 makes use of all the techniques you’ve learned about loops and conditionals.

Example 9-6. Iterating over a list of dictionaries in Jinja
{% for interface in interface_list %}
interface {{ interface.name }}
  description {{ interface.desc }}
  {% if interface.uplink %}
    switchport mode trunk
  {% else %}
    switchport access vlan {{ interface.vlan }}
    switchport mode access
  {% endif %}
{% endfor %}
TIP
When using Jinja to access data in a dictionary, you can use the traditional Python syntax dict['key'] or the shorthand form dict.key (as we’ve been showing). These two are identical. If you’re trying to access a key that doesn’t exist, a key error will be raised. However, you can also use the get() method in Jinja if it’s an optional key or if you want to return another value if the key doesn’t exist—for example, dict.get(key, 'UNKNOWN').

As mentioned previously, it’s bad form to embed data into your Python applications (see the interfaces list of dictionaries from the previous example). Instead, let’s place that data into its own YAML file and rewrite our application to import this data before using it to render the template. This is a good practice, because it allows someone with no Python experience to edit the network configuration by simply changing this simple YAML file.

Example 9-7 shows a YAML file that is written to be identical to your interfaces list in Example 9-6.

Example 9-7. YAML file for switchport configurations
---
- name: GigabitEthernet0/1
  desc: uplink port
  uplink: true
- name: GigabitEthernet0/2
  desc: Server port number one
  vlan: 10
- name: GigabitEthernet0/3
  desc: Server port number two
  vlan: 10
As we explored in Chapter 8, importing a YAML file in Python is easy. As a refresher, here’s our full Python application, but instead of the static, embedded list of dictionaries, you’re simply importing a YAML file to get that data:

from jinja2 import Environment, FileSystemLoader
import yaml

ENV = Environment(loader=FileSystemLoader('.'))

template = ENV.get_template("template.jinja")

with open("data.yml") as f:
    interfaces = yaml.safe_load(f)
    print(template.render(interface_list=interfaces))
You can reuse the same template from Example 9-6 and achieve the same result—but this time, the data that you’re using to populate your template comes from an external YAML file, which is a bit easier to maintain. The Python file now contains only the logic of pulling in data and rendering a template. This makes for a more maintainable template-rendering system.

This covers the basics of loops and conditionals. In this section, we’ve explored only a portion of what’s possible. Explore these concepts on your own and apply them to your own use cases.

Using Jinja Filters
The data structures (like lists and dictionaries) and basic logical constructs (like loops and conditionals) presented so far are useful, but sometimes something a little more advanced is warranted. What if you want to convert a given bit of input to uppercase before it’s rendered out to the template? Or perhaps you want to reverse the characters in a string? Tasks like these can be prohibitively difficult or even impossible with what we’ve covered in this chapter up to this point.

Jinja filters provide an extremely powerful answer to these and many more problems. You can think of them like little modular functions you can call upon at any point in your templates. They generally work by taking in data or text as input (such as a variable you are passing in to the template environment), performing some kind of specialized task on that data, and then outputting the result.

The syntax for using filters in a Jinja template is simple. In the same way that you can pipe output from a command into another command in a terminal shell on a Linux distribution, you can take the result of a Jinja statement and pipe it into a filter. The resulting text will make it into the rendered output of your template.

Example 9-8 uses the familiar double-brace syntax to specify that you want to use a variable named hostname. However, before the contents of hostname are rendered into the final output for this template, they are first piped (using |) into the filter called myfilter. The output from that filter is then used to populate that portion of the template.

Example 9-8. Basic filter syntax in Jinja
{{ hostname|myfilter }}
This is a bit contrived, so let’s move on to practical examples of Jinja filters in action.

Using the upper Jinja filter
Let’s take the preceding template and use a built-in filter called upper to capitalize the descriptions for each interface configuration:

{% for interface in interface_list %}
interface {{ interface.name }}
  description {{ interface.desc|upper }}
  {% if interface.uplink %}
    switchport mode trunk
  {% else %}
    switchport access vlan {{ interface.vlan }}
    switchport mode access
  {% endif %}
{% endfor %}
This filter, built into the Jinja2 library for Python, capitalizes the text piped to it.

Chaining Jinja filters
You can also chain filters, in much the same way that you might pipe commands together in Linux. Let’s use the reverse filter to take your capitalized text and print it backward:

{% for interface in interface_list %}
interface {{ interface.name }}
  description {{ interface.desc|upper|reverse }}
  {% if interface.uplink %}
    switchport mode trunk
  {% else %}
    switchport access vlan {{ interface.vlan }}
    switchport mode access
  {% endif %}
{% endfor %}
This results in the following output:

interface GigabitEthernet0/1
  description TROP KNILPU
  switchport mode trunk

interface GigabitEthernet0/2
  description ENO REBMUN TROP REVRES
  switchport access vlan 10
  switchport mode access

interface GigabitEthernet0/3
  description OWT REBMUN TROP REVRES
  switchport access vlan 10
  switchport mode access
To recap, our original description for GigabitEthernet0/1 was first uplink port, and then it was UPLINK PORT because of the upper filter, and then the reverse filter changed it to TROP KNILPU, before the final result was printed into the template instance.

Creating custom Jinja filters
This is all great, and tons of other great built-in filters are documented within the Jinja specification. But what if you want to create your own filter? Perhaps there’s something specific to network automation that you would like to perform in your own custom filter that doesn’t come with the Jinja2 library?

Fortunately, the library allows for this. Example 9-9 shows a full Python script defining a new function, get_interface_speed(). This function is simple—it looks for certain keywords like gigabit or fast in a provided string argument and returns the current Mbps value. It also loads all your template data from a YAML file, as shown in previous examples.

Example 9-9. Full Python script with custom Jinja filter
from jinja2 import Environment, FileSystemLoader         1
import yaml                                              1

ENV = Environment(loader=FileSystemLoader('.'))          2


def get_interface_speed(interface_name):
    """ get_interface_speed returns the default Mbps value for a given
        network interface by looking for certain keywords in the name
    """

    if 'gigabit' in interface_name.lower():
        return 1000
    if 'fast' in interface_name.lower():
        return 100

ENV.filters['get_interface_speed'] = get_interface_speed 3
template = ENV.get_template("template.jinja")

with open("data.yml") as f:                              4
    interfaces = yaml.safe_load(f)
    print(template.render(interface_list=interfaces))
1
Imports the Jinja2 and PyYAML libraries.

2
Declares the template environment.

3
Filters are added to the ENV object after declaration. We’re actually passing in the get_interface_speed function and not running it—the template engine will execute this function when we call template.render().

4
We load the YAML file and pass it in to the template when rendering it.

With a slight modification to our template, as shown in Example 9-10, we can leverage this filter by passing interface.name into the get_interface_speed filter. The resulting output will be whatever integer the function decides to return. Since all interface names are Gigabit Ethernet, the speed is set to 1000.

Example 9-10. Updated Jinja template leveraging the custom filter
{% for interface in interface_list %}
interface {{ interface.name }}
  description {{ interface.desc|upper|reverse }}
  {% if interface.uplink %}
  switchport mode trunk
  {% else %}
  switchport access vlan {{ interface.vlan }}
  switchport mode access
  {% endif %}
  speed {{ interface.name|get_interface_speed }}
{% endfor %}
TIP
You don’t always have to create your own function to pass in as a Jinja filter. Plenty of libraries provide helpful functions that you can import and simply pass right into the Jinja2 environment.

Template Inheritance in Jinja
As you create bigger, more capable templates for your network configuration, you may want to be able to break templates into smaller, more specialized pieces. It’s common to have a template for VLAN configuration, one for interfaces, and maybe another for a routing protocol. This kind of organizational tool, while optional, can allow for much more flexibility. The question is, how do you link these templates together in a meaningful way to form a full configuration?

Jinja allows you to perform inheritance in a template file, which is a handy solution to this problem. For instance, you may have a vlans.jinja file that contains only the VLAN configuration, and you can inherit this file to produce a VLAN configuration in another template file. You might be writing a template for interface configuration, and you want to also produce a VLAN configuration from another template. Example 9-11 shows how this is done by using the include statement.

Example 9-11. Using include for Jinja template inheritance
{% include 'vlans.jinja' %}

{% for name, desc in interface_dict.items() %}
interface {{ name }}
  description {{ desc }}
{% endfor %}
This renders vlans.jinja and inserts the resulting text into the rendered output for the template that includes it. Using the include statement, template writers can compose switch configurations made up of modular parts. This is great for keeping multiple template files organized.

Another inheritance tool in Jinja is the block statement. This is a powerful but more complicated method of performing inheritance, as it mimics object inheritance in more formal languages like Python. Using blocks, you can specify portions of your template that may be overridden by a child template, if present. If a child template is not present, this portion of the rendered output will still contain some default text.

Example 9-12 shows how blocks are declared within a parent template.

Example 9-12. Jinja block definition
{% for interface in interface_list %}
interface {{ interface.name }}
  description {{ interface.desc }}
{% endfor %}
!
{% block http %}
  no ip http server
  no ip http secure-server
{% endblock %}
Call this template no-http.jinja, indicating that you’d like to normally turn off the embedded HTTP server in your switch. However, you can use blocks to provide greater flexibility here. You can create a child template called yes-http.jinja that is designed to override this block and output the configuration that enables the HTTP server if that’s what you want; see Example 9-13.

Example 9-13. Jinja block reference
{% extends "no-http.jinja" %}
{% block http %}
  ip http server
  ip http secure-server
{% endblock %}
This allows you to enable the HTTP server simply by rendering the child template. The first line in Example 9-13 extends the parent template no-http.jinja so all of the interface configurations will still be present in the rendered output. However, because you’ve rendered the child template, the http block of the child overrides that of the parent. Using blocks in this way is useful for portions of the configuration that may need to change but aren’t properly served by traditional variable substitution.

The Jinja documentation on template inheritance goes into much more detail and is a great resource to keep bookmarked.

Variable Creation in Jinja
Jinja allows you to create variables within a template by using the set statement. A common use case is variable shortening. Sometimes you have to go through several nested dictionaries or Python objects to get what you want, and you may want to reuse this value several times in your template. Rather than repeat a long string of properties or keys, use the set statement to represent a particular value with a much smaller name:

{% set int_desc = sw01.config.interfaces['ge0/1']['description'] %}
{{ int_desc }}
While Jinja is clearly the dominant templating tool in network automation, it isn’t the only one. At times another tool is more appropriate, and we cover a few of them in the following sections.

Extensible Stylesheet Language Transformations
As discussed in Chapter 8, XML is an extremely popular data format that’s typically well supported by major NOSs for automation purposes and more. XML enjoys a rich ecosystem of tools surrounding it.

One of these is a fairly robust templating format called Extensible Stylesheet Language Transformations (XSLT). XSLT is typically used for applying transformations to XML data, such as to convert that data into XHTML or other XML documents. However, like Jinja, it can be used to create templates for any arbitrary document format, and given its affinity to the XML ecosystem, is an important option to have in your repertoire if you’re going to be working with XML frequently.

NOTE
As with XML itself, the XSLT specification is defined by W3C.

Let’s look at a practical example of populating an XSLT template with meaningful data in Python so that a resulting document can be achieved. The first thing you need is raw data to populate your template. This XML document will suffice:

<?xml version="1.0" encoding="UTF-8"?>
<authors>
    <author>
        <firstName>Christian</firstName>
        <lastName>Adell</lastName>
    </author>
    <author>
        <firstName>Scott</firstName>
        <lastName>Lowe</lastName>
    </author>
    <author>
        <firstName>Matt</firstName>
        <lastName>Oswalt</lastName>
    </author>
</authors>
This amounts to a list of authors, each with <firstName> and <lastName> elements. The goal is to use this data to generate an HTML table that displays these authors, via an XSLT document. An XSLT template to perform this task might look like Example 9-14.

Example 9-14. XSLT template
<xsl:stylesheet xmlns:xsl="http://www.w3.org/1999/XSL/Transform" version="1.0">
<xsl:output indent="yes"/>
<xsl:template match="/">
  <html>
  <body>
  <h2>Authors</h2>
    <table border="1">
      <tr bgcolor="#9acd32">
        <th style="text-align:left">First Name</th>
        <th style="text-align:left">Last Name</th>
      </tr>
      <xsl:for-each select="authors/author">
      <tr>
        <td><xsl:value-of select="firstName"/></td>
        <td><xsl:value-of select="lastName"/></td>
      </tr>
      </xsl:for-each>
    </table>
  </body>
  </html>
</xsl:template>
</xsl:stylesheet>
A few notes on this:

A basic for-each construct is embedded in what otherwise looks like valid HTML. This is a standard practice in template language—the static text remains static, and little bits of logic are placed where needed.

The for-each statement uses a coordinate argument (listed as authors/author) to state exactly which part of your XML document contains the data you wish to use. This is called XPath, a syntax used within XML documents and tools to specify a location within an XML tree.

You use the value-of statement to dynamically insert a value (like a variable in a Python program) as text from your XML data.

Assuming your XSLT template is saved as template.xslt and your data file as data.xml, you can return to your trusty Python interpreter to combine these two pieces and come up with the resulting HTML output:

from lxml import etree
xslRoot = etree.fromstring(bytes(open("template.xslt").read(), encoding='utf8'))
transform = etree.XSLT(xslRoot)

xmlRoot = etree.fromstring(bytes(open("data.xml").read(), encoding='utf8'))
transRoot = transform(xmlRoot)
This produces a valid HTML table, as shown in Figure 9-2.

npa2 0902
Figure 9-2. HTML table produced by XSLT
XSLT also provides additional logic statements:

<if>
Output the given element(s) only if a certain condition is met.

<sort>
Sort elements before writing them as output.

<choose>
A more advanced version of the if statement (allows else if or else style of logic).

You can take this example even further and use this concept to create a network configuration template using configuration data defined in XML; see Example 9-15.

Example 9-15. XML interface data
<?xml version="1.0" encoding="UTF-8"?>
<interfaces>
    <interface>
        <name>GigabitEthernet0/0</name>
        <ipv4addr>192.168.0.1 255.255.255.0</ipv4addr>
    </interface>
    <interface>
        <name>GigabitEthernet0/1</name>
        <ipv4addr>172.16.31.1 255.255.255.0</ipv4addr>
    </interface>
    <interface>
        <name>GigabitEthernet0/2</name>
        <ipv4addr>10.3.2.1 255.255.254.0</ipv4addr>
    </interface>
</interfaces>
Then, you can create an XSLT template to take this data and render it out into a new document that contains your valid network configuration, as shown in Example 9-16.

Example 9-16. XSLT template for router config
<xsl:stylesheet version="1.0" xmlns:xsl="http://www.w3.org/1999/XSL/Transform">
    <xsl:template match="/">
        <xsl:for-each select="interfaces/interface">
interface <xsl:value-of select="name"/>
    ip address <xsl:value-of select="ipv4addr"/>
        </xsl:for-each>
    </xsl:template>
</xsl:stylesheet>
With these XML and XSLT documents, you can get a rudimentary router configuration in the same way you generated an HTML page:

interface GigabitEthernet0/0
    ip address 192.168.0.1 255.255.255.0
interface GigabitEthernet0/1
    ip address 172.16.31.1 255.255.255.0
interface GigabitEthernet0/2
    ip address 10.3.2.1 255.255.254.0
As you can see, it’s possible to produce a network configuration by using XSLT. However, it is admittedly a bit cumbersome. You will likely find Jinja a much more useful templating language for creating network configurations, as it has a lot of features that are conducive to network automation.

Go Templates
The Go programming language has a package for working with templates in the standard library, importable as text/template. This package implements a template language that is tightly integrated with the Go programming language. So, if you’re working with Go, this package is a strong choice for building and rendering templates.

While the template language implemented by this package is specific to Go, and therefore includes some familiar Go-esque keywords and syntax, it does share some familiar concepts with other template languages like Jinja. For instance, the actual text replacement syntax in Go templates also uses double braces: {{ }}.

To get started, you’ll create a simple template using a string literal defined in your Go program, and print the results to the terminal; see Example 9-17.

Example 9-17. Simple Go template
package main

import (
  "os"
  "text/template"
)

func main() {
  // We can create an inline template using the Parse() method of the
  // template.Template type
  //
  // Note that simpleTemplate is just an arbitrary name chosen for this example.
  tmpl, err := template.New("simpleTemplate").Parse(`{{ "foobar" | print }}`)
  if err != nil {
    panic(err)
  }

  // We can render the template with Execute, passing in os.Stdout
  // as the first parameter, so we can see the results in our terminal
  err = tmpl.Execute(os.Stdout, nil)
  if err != nil {
    panic(err)
  }
}
NOTE
Example 9-17 shows a full working Go template that can be run with go run, but the following examples are simplified to show only the most relevant code. Full working versions of the following examples can be found at https://github.com/oreilly-npa-book/examples/tree/v2/ch09-templates.

You’re not even doing any variable substitution here. Go templates allow you to define string literals within the template itself, thus the escaped double quotes. So, this example simply results in the text foobar.

Also note that you don’t have to send the results to stdout. Anything that satisfies io.Writer can be used as the first parameter to Execute(), including a file opened via os.Create(), as shown in Example 9-18.

Example 9-18. Rendering a template to a file
tmpl, err := template.New("simpleTemplate").Parse(`{{ "foobar" | print }}`) 1
if err != nil {
  panic(err)
}

file, err := os.Create("./output.txt")                                      2
if err != nil {
  panic(err)
}

err = tmpl.Execute(file, nil)                                               3
if err != nil {
  panic(err)
}
1
Creates an inline template by using the Parse() method of the template​.Tem⁠plate type.

2
Creates the file to contain the output.

3
Anything that satisfies io.Writer can be passed as the first parameter to Execute(), which includes the file returned by os.Create.

To keep the remaining examples simple, you’ll continue to write to stdout.

Now, on to some more practical examples. As with Jinja, the chief use case for Go templates in network automation is to inject some kind of structured data into a predefined output format. This allows you to define the text that should more or less remain consistent within the template, and pass data retrieved from a user parameter or API call to “fill in” the parts of the template that require it.

For instance, you may have a Switch struct with the fields Hostname and InterfaceCount that you wish to use to pass data into your template; see Example 9-19.

Example 9-19. Passing data into a template via a struct
type Switch struct {                1
  Hostname       string
  InterfaceCount uint
}

sw01 := Switch{"sw01", 48}          2

tmpl, err := template.New(
    "switchTemplate").Parse("Device {{.Hostname}} has {{.InterfaceCount}} interfaces\n"
)                                   3
if err != nil {
  panic(err)
}

err = tmpl.Execute(os.Stdout, sw01) 4
if err != nil {
  panic(err)
}
1
First, we must define the struct that represents a switch.

2
Then we instantiate this struct into a variable sw01.

3
We can refer to the fields of the struct by using the .field name syntax.

4
Our input struct is passed in as the second parameter to Execute().

You may have noticed you were able to refer to the Hostname field of your struct via a leading dot (.). In the package documentation, this is literally referred to as dot and is a reference to whatever type is passed in to your template context. In this case, dot is a simple struct with fields you could refer to by name, but this doesn’t have to always be the case, as you’ll see in the next few examples.

This prints the following:

Device sw01 has 48 interfaces
Instead of a single struct, you can pass in a slice of multiple Switch elements and iterate over them using the range keyword:

type Switch struct {
  Hostname       string
  InterfaceCount uint
}

// switches is a slice that represents all our Switch instances
switches := []Switch{
  {"sw01", 48},
  {"sw02", 24},
  {"sw03", 48},
}

// As with Jinja, it's often better to define templates in a multiline string
// like tmplStr below, or in a separate file that is read in before rendering.
tmplStr := `
{{range $i, $switch := .}}
Device {{$switch.Hostname}} has {{$switch.InterfaceCount}} interfaces
{{end}}
`
In this case, dot represents your slice, and iterating over this slice in the template is done via range $i, $switch := . This is similar to iterating over a slice natively in Go, as the slice is to the right of the short variable declaration operator :=, and i and switch represent the index value of each element and a copy of the element itself, respectively.

Go templates support conditionals as well. You can optionally print a line based on the value of the Enabled boolean field:

type Switch struct {
  Hostname       string
  InterfaceCount uint
  Enabled        bool
}

// switches is a slice that represents all our Switch instances
switches := []Switch{
  {"sw01", 48, true},
  {"sw02", 24, true},
  {"sw03", 48, false},
}

tmplStr := `
{{range $i, $switch := .}}
{{if $switch.Enabled}}
Device {{$switch.Hostname}} has {{$switch.InterfaceCount}} interfaces
{{end}}
{{end}}
`
Your struct thus far has been fairly simple, and you might wonder if you can simplify it to a map and iterate over that. You can!

switches := map[string]int{
  "sw01": 48,
  "sw02": 24,
  "sw03": 48,
}

// Since we're now iterating over a map, the two created variables after the
// range keyword represent each key/value pair
tmplStr := `
{{range $hostname, $ifCount := .}}
Device {{$hostname}} has {{$ifCount}} interfaces
{{end}}
`
When discussing Jinja, we explored the use of filters to further manipulate data within the template and even defined our own custom filter. You can do the same in Go templates, albeit with a slightly different terminology. Go templates include a feature called pipelines, which can be a sequence of commands delimited by the pipe operator (|). Go has several built-in commands, but you’ll see how to define your own. While the terms differ, the syntax is extremely similar to Jinja.

Let’s go back to passing a Switch struct as the data parameter for your template, but instead of an integer field representing the number of interfaces, this struct can contain a slice of strings, so you can actually have the names of each interface. This will allow you to iterate over the interface names if you want to, but at the same time, you can also still print the total number of interfaces by passing that field through a pipeline, using the len command, as shown in Example 9-20.

Example 9-20. Go template pipelines
// A fairly small switch, indeed
sw01 := Switch{"sw01", []string{
  "ge-0/0/1",
  "ge-0/0/2",
  "ge-0/0/3",
  "ge-0/0/4",
}}

tmplStr := "Device {{.Hostname}} has {{ .Interfaces | len }} interfaces\n"
The built-in len command prints the length of the value passed to it, rather than the value itself. This results in the following:

Device sw01 has 4 interfaces
Other commands are available, such as slice for subslicing an argument, and and or for comparing two arguments for emptiness, and call for calling a function. However, you can create your own function and make it available to the template context on creation. For instance, you can create a function called IfParse() for parsing the speed and location details out of a string interface name:

type Interface struct {
  Speed int
  FPC   int
  PIC   int
  Port  int
}

func IfParse(ifaceStr string) Interface {

  iface := Interface{}

  ifSplit := strings.Split(ifaceStr, "-")

  speeds := map[string]int{
    "ge": 1,
    "xe": 10,
    "et": 40,
  }
  iface.Speed = speeds[ifSplit[0]]

  locSplit := strings.Split(ifSplit[1], "/")

  fpc, _ := strconv.Atoi(locSplit[0])
  iface.FPC = fpc

  pic, _ := strconv.Atoi(locSplit[1])
  iface.PIC = pic

  port, _ := strconv.Atoi(locSplit[2])
  iface.Port = port

  return iface
}
To use this function, you must create a mapping between the function itself and a name you can refer to in your template. The template package has a function called FuncMap() you can use to create this. Then, you must remember to pass in this function map by using the Funcs() method of the template object, as shown in Example 9-21.

Example 9-21. Using custom functions in Go templates
fmap := template.FuncMap{"ifparse": IfParse}                           1

sw01 := Switch{"sw01", []string{
  "ge-0/0/1",
  "ge-0/0/2",
  "ge-0/0/3",
  "ge-0/0/4",
}}

tmplStr := `
{{range $i, $interface := .Interfaces}}
{{with $loc := $interface | ifparse}}
Interface {{$interface}}   port {{$loc.Port}}
{{end}}
{{end}}
`                                                                      2

tmpl, err := template.New("switchTemplate").Funcs(fmap).Parse(tmplStr) 3
if err != nil {
  panic(err)
}
1
Creates a mapping of functions to names we can refer to in the template.

2
In this template, we’re first creating a new variable loc by pipelining $interface into the custom function ifparse. Since ifparse returns an Interface type, we can then refer to the fields of that struct directly!

3
Don’t forget to pass in the FuncMap by using the Funcs() method, as shown here.

This prints the following:

Interface ge-0/0/1   port 1
Interface ge-0/0/2   port 2
Interface ge-0/0/3   port 3
Interface ge-0/0/4   port 4
There is much more to explore with Go templates, but we’ve covered the basics. Again, if you’re working with Go and already familiar with the basic concepts we previously discussed with Jinja, you could do worse than to skill up on Go templates.

Summary
You don’t have to work directly with programming languages like Go and Python in order to use templates, as we’ve done in this chapter. As we discuss in Chapter 12, tools like Ansible and Nornir include native support for template languages like Jinja2, allowing you to populate your templates with data available within those platforms.

Here are a few parting thoughts on using templates for network automation:

Keep the templates simple. Leveraging loops and conditionals to enhance your templates is fine, but don’t go overboard. Jinja isn’t as robust as a fully featured, general-purpose programming like Python, so keep the more advanced stuff out of the template.

Leverage template inheritance to reuse portions of configurations that don’t need to be duplicated.

Remember, syntax and data should be handled separately. For instance, keep VLAN IDs in their own data file (maybe YAML) and the CLI syntax to implement those VLANs in a dedicated template.

Templates are just text files, so you can and should use version control (i.e., Git) to track changes and collaborate with others. We discuss this in Chapter 11.

# Working with Network APIs
From Python, Go, and data formats to configuration templating with Jinja, we’ve explored key foundational technologies and skills that will make you a better network engineer. In this chapter, you’re going to put these skills to practical use and start to consume and communicate with various types of network device APIs to start automating your network.

As we introduced in Chapter 2, nowadays there are multiple options to interact with network platforms. Along with the traditional CLI and SNMP, we have new alternatives—from network-specific APIs (such as NETCONF, RESTCONF, and gNMI) to multipurpose APIs (such as HTTP-based ones or the Linux shell). Not every device supports all of these options, so understanding their capabilities will determine your automation options.

All interfaces are viable for automation, each one with its own pros and cons. The goal of this chapter is to introduce these APIs, showcasing how you can use them programmatically in Python and Go.

To best help you understand how to start interacting with networks programmatically, this chapter is organized into two sections:

Understanding network APIs
We examine the architecture and foundation of APIs, including RESTful and non-RESTful HTTP-based APIs, NETCONF, RESTCONF, and gRPC/gNMI. In each case, we introduce common tools used for testing and show how to use each one.

Automating using network APIs
We introduce some popular Python and Go libraries that allow you to start creating applications to interact with your network. We’ll look at the Python Requests and Go HTTP libraries for consuming HTTP-based APIs (including RESTCONF), the Python ncclient for interacting with NETCONF devices, the Go gNMIc for interacting with the gNMI interface, and the Python Netmiko library for automating devices over SSH.

As you read this chapter, keep in mind one thing: this chapter is not a comprehensive guide on any particular API and should not serve as API documentation. We provide examples using different vendor implementations of a given API, as it’s common to be working in a multivendor environment. It’s also important to see the common patterns and unique contrasts among APIs.

Understanding Network APIs
Our focus is on four of the most common types of APIs you’ll find on network devices: HTTP-based APIs, NETCONF, RESTCONF, and gRPC/gNMI. We’re going to start by looking at foundational concepts for each type of API; once we review them, we’ll explore the consumption of these APIs with hands-on examples using multiple vendors.

NOTE
For each network API type, we have used one or two network platforms. This doesn’t imply that each API is the only interface a platform supports. Actually, each platform usually supports multiple interfaces, but for illustrating multiple vendors and interfaces, we have used an arbitrary mapping to show diversity, without extra considerations.

As we start our journey of consuming and interacting with network APIs, in each API subsection, our focus is just like the focus we’ve had thus far throughout the book—​on vendor-neutral tools and libraries. More specifically, we are going to look at tools such as cURL for working with HTTP-based APIs (RESTCONF included), NETCONF over SSH for working with NETCONF APIs, and gNMIc to interact with the gNMI interface.

It’s important to note that this section is about exploring network APIs in that we showcase how to get started using and testing network APIs without writing any code. We want you to understand the concepts from each particular API type before putting them to use in the next section. This section is not about the tools and techniques you would use for automating production networks. Those types of tools and libraries are covered in “Using Network APIs for Automation”.

Let’s get started by diving into HTTP-based APIs.

Getting Familiar with HTTP-Based APIs
HTTP-based APIs are not exclusively used for network management. They are one of the most common interprocess connection types; thus most of the concepts introduced in this section apply to general use cases. Within the context of network automation, you will learn how to use APIs to manage network services using HTTP APIs as the management interface. For instance, HTTP APIs are used in Chapter 12 to provision dynamic network infrastructure via Terraform providers. In the same chapter, HTTP APIs are used to fetch data from a source of truth (SoT) containing the network device inventory and to create a dynamic inventory for Nornir.

You should understand two types of HTTP-based APIs in the context of network APIs: RESTful HTTP-based APIs and non-RESTful HTTP-based APIs. To better understand them and what the term RESTful means, we are going to start by examining RESTful APIs. Once you understand RESTful architecture and principles, we’ll move on and compare them with non-RESTful HTTP-based APIs.

Understanding RESTful APIs
RESTful APIs are becoming more popular and more commonly used in the networking industry, although they’ve been around since the early 2000s. Most of the APIs that exist today within network infrastructure are HTTP-based RESTful APIs. Therefore, when you hear about a RESTful API on a network device or SDN controller, it is an API that will be communicating between a client and a server.

The client is an application such as a Python script or web UI application, and the server is the network device or controller. Moreover, since HTTP is being used as transport, you’ll perform some operations using URLs just as you do already as you browse the internet. Thus, if you understand that when you’re browsing a website, HTTP GETs are performed, and when you’re filling out a web form and clicking Submit, an HTTP POST is performed, you already understand the basics of working with RESTful APIs.

Let’s look at examples of retrieving data from a website and retrieving data from a network device via a RESTful API. In both instances, an HTTP GET request is sent to the web server (see Figure 10-1).

In Figure 10-1, one of the primary differences is the data that is sent to and from the web server. When browsing the internet, you receive HTML data that your browser will interpret so that it can properly display the website. On the other hand, when issuing an HTTP GET request to a web server that is exposing a RESTful API (remember, it’s exposing it via a URL), you receive data back that is mostly encoded using JSON or XML. This is where you’ll use what we reviewed in Chapter 8. Since you receive data back in JSON/XML, the client application must understand how to interpret JSON and/or XML. Let’s continue with the overview, so you have a more complete picture before we start to explore the use of RESTful HTTP APIs.

![Alt text](image.png)\
\
Figure 10-1. Understanding REST by looking at HTTP GET responses (npa2 1001)
Let’s take our high-level overview one step further and look at the origins of RESTful APIs. The birth and structure of modern web-based RESTful APIs came from a PhD dissertation by Roy Fielding in 2000. In “Architectural Styles and the Design of Network-based Software Architectures,” he defined the intricate detail of working with networked systems on the internet that use the architecture defined as REST.

An interface must conform to six architectural constraints in order to be considered RESTful. For the purposes of this chapter, we’ll look at three:

Client-server
This is a requirement to improve the usability of systems while simplifying the server requirements. Having a client-server architecture allows for the portability and changeability of client applications without the server components being changed. This means you could have different API clients (web UI, CLI) that consume the same server resources (backend API).

Stateless
The communication between the client and server must be stateless. Clients that use stateless forms of communication must send all data required for the server to understand and perform the requested operation in a single request. This is in contrast to interfaces such as SSH, which have a persistent connection between a client and a server.

Uniform interface
Individual resources in scope within an API call are identified in HTTP request messages. For example, in RESTful HTTP-based systems, the URL used references a particular resource. In the context of networking, the resource maps to a network device construct such as a hostname, interface, routing protocol configuration, or any other resource that exists on the device. The uniform interface also states that the client should have enough information about a resource to create, modify, or delete a resource.

These are just three of the six core constraints of the REST architecture, but you likely can already see the similarity between RESTful systems and how you consume the internet through web browsing on a daily basis. Keep in mind that HTTP is the primary means of implementing RESTful APIs, although the transport type could, in theory, be something else. To really understand RESTful APIs, then, you must also understand the basics of HTTP.

Understanding HTTP request types
While every RESTful API you look at is an HTTP-based API, you will eventually look at HTTP-based APIs that do not adhere to the principles of REST and therefore are not RESTful. In either case, the APIs require an understanding of HTTP. Because these APIs are using HTTP as transport, you’re going to be working with the same HTTP request types and response codes that are used on the internet already.

Common HTTP request types include GET, POST, PATCH, PUT, and DELETE. As you can imagine, GET requests are used to request data from the server, DELETE requests are used to delete a resource on the server, and the three Ps (POST, PATCH, PUT) are used to make a change on the server. In Table 10-1, we list each method’s definition along with its meaning in the context of networking.

Table 10-1. HTTP request
Request type	Description	In networking context
GET

Retrieves a specified resource

Obtaining configuration or operational data

PUT

Creates or replaces a resource

Making a configuration change

PATCH

Creates or updates a resource object

Making a configuration change

POST

Creates a resource object

Making a configuration change

DELETE

Deletes a specified resource

Removing a particular configuration

Understanding HTTP response codes
Just as the request types are the same if you’re using a web browser on the internet or using a RESTful API, the same is true for response codes.

Ever see a 401 Unauthorized message when you were trying to log in to a website and used invalid credentials? Well, you would receive the same response code if you were trying to log in to a system using a RESTful API and you sent the wrong credentials. The same is true for successful messages or if the server has an error of its own. Table 10-2 lists the common types of response codes you see when working with HTTP-based APIs. This list is not exclusive; others exist too.

Table 10-2. HTTP response codes
Response code	Description
1XX

Informational

2XX

Successful

3XX

Redirect

4XX

Client error

5XX

Server error

Remember, the response code types for HTTP-based APIs are no different from standard HTTP response codes. We are merely providing a list of the types and will leave it as an exercise for you to learn about individual responses.

Exploring HTTP-based APIs with cURL
cURL is a command-line tool for working with URLs. From the Linux command line, you can send HTTP requests by using the cURL program. While cURL uses URLs, it can communicate to servers using protocols besides HTTP, including FTP, SFTP, TFTP, Telnet, and many more. You can use the command man curl from the Linux command line to get an in-depth look at all the various options cURL supports.

NOTE
cURL isn’t limited to Linux. It’s available in multiple OSs such as macOS and Windows. For installation instructions, check https://curl.se/docs/install.html.

Multiple alternatives to cURL are available, either as command-line tools or GUIs, but all share the same concepts. Once you understand the basic ideas, you can apply them to the other tools. However, using a user-intuitive web GUI frontend, such as Postman, could make it much easier to learn and test HTTP APIs. These GUI tools put the focus on using the API without worrying about writing code. You’ll see an example shortly in Figure 10-2 to help you understand the look and feel.

We start our exploration of HTTP APIs by using cURL with the Cisco Meraki RESTful API (Meraki’s API documentation is available at https://oreil.ly/zHOWy). Cisco Meraki is a cloud networking controller that helps to illustrate how to interact with this type of network infrastructure. Many modern NOSs also offer REST APIs, usually (but not limited to) implementing the RESTCONF interfaces, covered in “Using RESTCONF”.

Using the HTTP GET method to retrieve information
As we’re just getting started with RESTful APIs, we’ll begin with a simple HTTP GET request to retrieve all organizations from the API, targeting the URL https://api.meraki.com/api/v1/organizations.

NOTE
Organization is an abstract concept from Cisco Meraki created to support multiple tenants for the same account. Each will contain different network resources.

In Example 10-1, we use a cURL statement to call the Cisco Meraki URL and retrieve a list of all the organizations.

NOTE
Full versions of the code examples in this chapter can be found in the book’s GitHub repo at https://github.com/oreilly-npa-book/examples/tree/v2/ch10-apis.

Example 10-1. Retrieving Meraki organizations with cURL
$ curl 'https://api.meraki.com/api/v1/organizations' \                    1 2
  -H 'X-Cisco-Meraki-API-Key: 6bec40cf957de430a6f1f2baa056b99a4fac9ea0' \ 3
  -L                                                                      4

# response omitted
1
The URL is generic; it is shared by all Cisco Meraki customers. It’s offered as IaaS, covered in Chapter 4.

2
We have not defined any HTTP operation. Nevertheless, by default, cURL performs a GET operation. This behavior can be modified using the -X flag, as in next examples. You can see all the available cURL customizations via the command documentation: man cURL.

3
The -H argument, or --header, is used to include HTTP headers in the HTTP request. HTTP headers are key-value pairs used to pass metadata to the server, useful for things like authentication.

4
The -L flag, or --location, allows the client to follow any redirects issued by the server.

NOTE
The Cisco Meraki API token used in Example 10-1 has been taken from Cisco Developer Hub. You can use the same one if it is still active. If not, a new one will likely be available on this website, for API exploration.

Also, note that in the previous URL, the base URL path contains /v1/. This is an arbitrary way to indicate the targeted version of this API. As with any other application, the API can evolve over time, adding, changing, and removing resources. Using API versioning is a common pattern to offer a predictable behavior to consumer applications, without facing breaking changes (i.e., accessing a path that has been removed). In other APIs, the version may be specified via the api_version query parameter. This way, the URL path is not modified, and only the query parameter is appended; here’s an example:

$ curl https://my_application.com/api/my_path?api_version=1.3
The omitted output from the cURL statement in Example 10-1 is an output word wrapped on the terminal, which is hard to read. Alternatively, as shown in Example 10-2, you can pipe the response to python3 -m json.tool to pretty-print the response object, making it much more human-readable.

Example 10-2. Using the Python json.tool to render a JSON response
$ curl 'https://api.meraki.com/api/v1/organizations' \
  -H 'X-Cisco-Meraki-API-Key: 6bec40cf957de430a6f1f2baa056b99a4fac9ea0' \
  -L \
  | python3 -m json.tool

[
  {
    "id": "573083052582915028",
    "name": "Next Meraki Org",
    "url": "https://n18.meraki.com/o/PoiDucs/manage/organization/overview",
    "api": {
      "enabled": true
    },
    "licensing": {
      "model": "co-term"
    },
    "cloud": {
      "region": {
          "name": "North America"
      }
    }
  },
  # other organizations omitted for brevity
]
The object retrieved is a JSON object, which converts to a list in Python because it begins and ends with square brackets. Each item in the list is a dictionary, representing an organization and all its attributes. The response media type (in this case, JSON), can be influenced by the Accept header (expressing the client wish), but in this case, it has no impact because JSON is the only media type supported by this API.

To compare the user experience from cURL to Postman, Figure 10-2 shows an equivalent HTTP GET request done via the Postman GUI.

![Alt text](image-1.png)\
\
Figure 10-2. Postman GET request (npa2 1002)
WARNING
The UI you get from Postman may differ from the one in this book. UIs are evolving with the product, so it’s likely to change over time. However, the concepts remain the same.

In Figure 10-2, you can appreciate the same request and output from Example 10-1, but in a better visual presentation. Following the same pattern, you could reproduce all the examples in this section in Postman.

TIP
Postman allows you to create and publish Postman Collections as common API examples to be reused (using variables for customization). These collections can serve as a good reference for common operations on APIs. As an example, Nick Russo maintains an interesting collection for several network APIs at https://oreil.ly/QBJmd.

Commonly, behind a REST API, different resources are related. From Cisco Meraki documentation, we know that each organization can contain networks. So, similar to the previous API endpoint for organizations, there is one for networks: /api/v1/organizations/{organizationId}/networks. Between the curly braces, the organizationID must be replaced by the actual organization identifier.

This identifier is the most relevant attribute for each organization that you retrieved in Example 10-2. It is represented by the id key in each dictionary. Taking this ID, you can continue exploring the nested networks that belong to each organization. It’s interesting to notice the nested nature of this API, so you can’t list all the networks directly, but use the organizations they belong to as a reference. In Example 10-3, the organization ID is used to retrieve the networks that belong to it.

Example 10-3. Retrieving Meraki networks with cURL
$ curl 'https://api.meraki.com/api/v1/organizations/573083052582915028/networks' \ 1
  -H 'X-Cisco-Meraki-API-Key: 6bec40cf957de430a6f1f2baa056b99a4fac9ea0' \
  -L \
  | python3 -m json.tool

[
  {                                                                                2
    "id": "L_573083052582989052",                                                  3
    "organizationId": "573083052582915028",
    "name": "Long Island Office",
    "productTypes": [
      "appliance",
      "camera",
      "switch"
    ],
    "timeZone": "America/Los_Angeles",
    "tags": [
      "tag1",
      "tag2"
    ],
    "enrollmentString": null,
    "url": "https://n18.meraki.com/Long-Island-Offi/n/kWaHAbs/manage/usage/list",
    "notes": "Combined network for Long Island Office",
    "isBoundToConfigTemplate": false
  },
  # other networks removed for brevity
]
1
The URL path contains the organization ID that limits the scope of the request to the networks belonging to that organization.

2
The response is a list of dictionaries, each representing a network. And in each dictionary, we can find each network’s attributes.

3
Similar to the previous organization’s example, the id key is used to uniquely identify the network.

Next, we continue exploring HTTP methods introduced in “Understanding HTTP request types”. In particular, we’ll start with a method you can use to modify resources on the API: the POST.

Using the HTTP POST method to create a new resource
In Example 10-3, you retrieved the networks belonging to a specific organization. Now, for the same organization, you want to create a new network (Example 10-4).

Example 10-4. Creating a Meraki network with cURL
curl -X POST 'https://api.meraki.com/api/v1/organizations/573083052582915028/networks' \ 1
  -H 'X-Cisco-Meraki-API-Key: 6bec40cf957de430a6f1f2baa056b99a4fac9ea0' \
  -L \
  -d '{"name": "my new automated network", "productTypes": ["switch"]}' \                2
  -H 'Content-Type: application/json' \                                                  3
  | python3 -m json.tool

{
    "id": "N_573083052583237701",
    "organizationId": "573083052582915028",
    "productTypes": [
        "switch"
    ],
    "url": "https://n18.meraki.com/my-new-automated/n/mQ9KWds/manage/usage/list",
    "name": "my new automated network",
    "timeZone": "America/Los_Angeles",
    "enrollmentString": null,
    "tags": [],
    "notes": null,
    "isBoundToConfigTemplate": false
}
1
The HTTP method to create new objects is POST, and it is specified with the -X flag. The POST method requires data.

2
With the -d flag, or --data, we pass a JSON object with the attributes of the new network in the form of a key-value pair (name and productTypes).

3
The Content-Type header is used to specify the data format. In this case, we are using JSON, but other formats are also supported (e.g., XML).

TIP
Learning how to construct a proper API request requires becoming familiar with API documentation. The API documentation (the API definition and specs) defines what a given URL must be, the HTTP request type, headers, and what the body needs to be for a successful API call. For instance, in the previous example, we passed the required attributes for the POST request, but we could have also passed optional attributes, such as timeZone or tags. All these attributes are defined in the API documentation. Additionally, performing GET requests offers some hints of the required attributes, as you can see in the output of the networks GET (in Example 10-3) and the data used for the POST.

Now that you understand the principles of REST and HTTP, it’s important to also take note of non-RESTFul HTTP-based APIs.

Understanding non-RESTful HTTP-based APIs
RESTful APIs are the most popular HTTP-based APIs, and other HTTP-based APIs are not compliant with REST principles. In the network industry, during the adoption of newer interfaces, such as RESTful ones, some APIs were built on top of CLIs, meaning that the API call actually sends a command to the device versus sending native structured data. Obviously, the preferred approach is to have any modern network platform’s CLI or web UI use the underlying API, but for legacy or preexisting systems that were built using commands, it is common to see the use of non-RESTful APIs, as it was easier to add an API this way rather than rearchitect the underlying system.

RESTful HTTP-based APIs and non-RESTful HTTP-based APIs have two major differences. RESTful APIs use particular verbs (e.g., GET, POST, PATCH, etc.) to dictate the type of change being requested of the target server. For example, in the context of networking, a configuration change would never occur if you’re doing an HTTP GET, since you’re simply retrieving data. However, systems that are HTTP based but do not follow RESTful principles could use the same HTTP verb for every API call. This means if you’re retrieving data or making a configuration change, all API calls could be using a POST request. Another common difference is that non-RESTful HTTP-based APIs always use the same URL and do not allow you to access a specific resource via a URL change. You can see both characteristics in Example 10-5.

Within the non-REST HTTP-based APIs, there is one popular methodology, the RPC, which was available before the REST APIs become popular. An RPC is a simple calling to a function in a remote system, with a data payload containing a method, and some other attributes. Depending on how the data is codified, we could talk about XML-RPC or JSON-RPC. This command-and-action approach makes it more performant, but also more obscure in terms of predictability.

NOTE
Both types, REST and RPC, can coexist on the same API server, in different parts of the API, leveraging their benefits for different use cases. We will present more RPC use cases in “Using NETCONF” and “Understanding gRPC”.

One example of a JSON-RPC API is the Arista eAPI. It offers an RPC endpoint (/command-api) to run CLI commands via the HTTP API. Example 10-5 uses cURL again to request the execution of the CLI commands providing the proper JSON payload.

Example 10-5. Running CLI commands via the Arista eAPI
$ curl --insecure \
  -H "Content-Type: application/json" \
  -X POST \                                              1
  -d '{"jsonrpc":"2.0", "method":"runCmds", "params":{ "version":1,
  "cmds":["show version"], "format":"text"}, "id":""}' \ 2
  https://ntc:ntc123@eos-spine1/command-api \            3
  | python3 -m json.tool

{
    "jsonrpc": "2.0",
    "id": "",
    "result": [
        {                                                4
            "output": " vEOS\nHardware version:    \nSerial number:       \n
            System MAC address:  5254. 0097.1b5e\n\nSoftware image version:
            4.22.4M\nArchitecture:           i686\nInternal build version:
            4.22.4M-15583082.4224M\nInternal build ID:
            08527907-ec51-458e-99dd-e3ad9c80cbbd\n\nUptime:
            15 weeks, 4 days, 13 hours and 22 minutes\nTotal memory:
            2014520 kB\nFree memory:            1335580 kB\n\n"
        }
    ]
}
1
We are using the POST method to retrieve data. As we’ve commented, in non-REST APIs, the method is not meaningful. Actually, with the same method, depending on the CLI commands passed, we could be retrieving the state or changing it.

2
We define the operation to be executed remotely with the method (runCmds), and the cmds parameter, which contains a list of all the commands to be executed.

3
The authentication parameters (ntc:ntc123@) are passed in the URL that is equivalent to the standard HTTP Authorization header.

4
The result key contains the output of the command executed—in this case, the output of the show version command without any formatting, simply the raw text (i.e., what we get via an SSH CLI access).

NOTE
Alternatives to RESTful APIs other than RPC have appeared. For instance, GraphQL was published in 2015 by Facebook. It defines a data query language that simplifies the way data is consumed, allowing clients to define the structure/filtering of the data required, which will be served by the server. This approach reduces the amount of data transferred but can impede caching of the results. GraphQL is especially useful for retrieving data from an SoT, collecting the relevant data from one object (including nested resources). We dig into this in more detail in Chapter 14.

As you start to use various types of HTTP-based APIs on network devices, keep in mind the following points:

HTTP APIs can use XML or JSON for data encoding, but the device may implement only one or the other. The API’s author determines what gets supported.

Tools such as cURL and Postman are helpful as you get started with APIs, but to write code to interact with HTTP APIs, you need a library that speaks HTTP, such as the Python Requests library or Go net/http package (covered in “Using Network APIs for Automation”).

Pay close attention to the HTTP verbs used when making configuration changes—using the wrong verb can have unintended consequences.

You need to use API documentation to understand how to construct a proper API request. You’ll need the URL, headers, HTTP method, and body.

Now that we’ve introduced HTTP-based APIs, let’s shift our focus and introduce the NETCONF API.

Using NETCONF
NETCONF is a network configuration management protocol, defined in RFC 6241 and designed from the ground up for configuration management and for retrieving configuration and operational state data from network devices. In this respect, NETCONF has a clear delineation between configuration and operational state; API requests are used to perform operations such as retrieving the configuration state, retrieving the operational state, and making configuration changes.

NOTE
We stated in the previous section that RESTful APIs aren’t new; they are merely new for network devices and SDN controllers. As we transition to looking at the NETCONF API, it’s worth noting that NETCONF is also not new. NETCONF has been around for nearly two decades. In fact, it’s an industry-standard protocol with its original RFC published in 2006. It’s even been on various network devices for years, although often as a limited API rarely being used.

One of the core attributes of NETCONF is its ability to utilize various configuration data stores. Most network engineers are familiar with running configurations and startup configurations. These are thought of as two configuration files, but they are two configuration data stores in the context of NETCONF.

NETCONF implementations often tend to use a third data store called a candidate configuration. The candidate configuration data store holds configuration objects (CLI commands if you’re using CLI for configuration) that are not yet applied to the device. As an example, if you enter a configuration on a device that supports candidate configurations, they do not take action immediately. Instead, they are held in the candidate configuration and applied to the device only when a commit operation is performed. When the commit is executed, the candidate configuration is written to the running configuration.

Candidate configuration data stores have been around for years as originally defined in the NETCONF RFC almost two decades ago. One of the issues the industry has faced is having usable implementations of NETCONF that offered this functionality. However, not all implementations have been unused—​there have, in fact, been successful implementations. Juniper’s Junos OS has had a robust NETCONF implementation for years, along with the capability of a candidate configuration; more recently, other products from vendors such as Cisco, Huawei, and Nokia have adopted support for the candidate configuration data store.

TIP
Always check your hardware and software platforms, even if they are from the same vendor. The capabilities they support likely differ. The support of a candidate configuration is just one example.

We stated that with a candidate configuration, you enter various configurations, and they aren’t yet applied until a commit operation is performed. This leads us to another core attribute of NETCONF-enabled devices: configuration changes as a transaction. In our example, it means that all configuration objects (commands) are committed as a transaction. All commands succeed or are not applied. This is in contrast to the more common scenario of entering a series of commands and having a command somewhere in the middle fail, yielding a partial configuration.

Moreover, a unique feature of NETCONF (not available in any of the other interfaces) is the network-wide transaction. If a configuration change affects multiple devices—for example, provisioning a Layer 3 VPN (L3VPN) service end to end, updating the configuration of several network devices—all the changes need to succeed, or the whole operation is rolled back (which is called an abort phase).

The support of a candidate configuration and atomic transactions are just two features of NETCONF. Let’s take a deeper dive into the underlying NETCONF protocol stack.

Learning the NETCONF protocol stack
Now we’ll tackle the four basic layers in NETCONF: the content representation, the operation types, the messages, and the supported transport protocols, as outlined in Table 10-3. We are going to review each and show concrete examples of what they mean for the XML object being sent between the client and server.

Table 10-3. NETCONF protocol stack
Layer	Example
Content

XML representation of data models (YANG, XSD)

Operations

get-config, get, copy-config, lock, unlock, edit-config, delete-config, kill-session, close-session, commit, validate, ...

Messages

rpc, rpc-reply, hello

Transport

SSHv2, SOAP, TLS, BEEP

NOTE
NETCONF supports only XML for data encoding. On the other hand, remember that RESTful APIs have the ability to support JSON and/or XML.

Transport
NETCONF is commonly implemented using SSH as transport; it is its own SSH subsystem. While all of our examples use NETCONF over SSH, it is technically possible to implement NETCONF over SOAP, TLS, or any other protocol that meets the requirements of NETCONF.

A few of these requirements are as follows:

It must be a connection-oriented session, and thus there must be a consistent connection between a client and a server.

NETCONF sessions must provide a means for authentication, data integrity, confidentiality, and replay protection.

Although NETCONF can be implemented with other transport protocols, each implementation must support SSH at a minimum.

NOTE
With the popularity of RESTful APIs, instead of further developing NETCONF over other protocols, the brand-new RESTCONF interface was created to implement NETCONF functionalities using REST principles. We cover this approach in “Using RESTCONF”.

Messages
NETCONF messages are based on an RPC–based communication model, and each message is encoded in XML. Using an RPC-based model allows the XML messages to be used independent of the transport type. NETCONF supports three message types: <hello>, <rpc>, and <rpc-reply>. Viewing the actual XML-encoded object helps elucidate NETCONF, so let’s take a look at a NETCONF RPC request.

The <hello> message is sent by the NETCONF server when the connection is established, exposing its capabilities—the data models and the actions supported:

<hello xmlns="urn:ietf:params:xml:ns:netconf:base:1.0">
  <capabilities>
    <capability>urn:ietf:params:netconf:base:1.1</capability>
    <!-- rest of request as XML... -->
  </capabilities>
</hello>
NOTE
For a refresher on XML concepts, such as the XML namespace xmlns, see Chapter 8.

After the initial <hello>, the message types are always going to be <rpc> and <rpc-reply> and will always be the outermost XML tag in the encoded object:

<rpc message-id="101">
    <!-- rest of request as XML... -->
</rpc>
Every NETCONF <rpc> includes a required attribute called message-id. You can see this in the preceding example. It’s an arbitrary string the client sends to the server. The server reuses this ID in the response header so the client knows which message the server is responding to.

The other message type is <rpc-reply>. The NETCONF server responds with the message-id and any other attributes received from the client (e.g., XML namespaces):

<rpc-reply message-id="101" xmlns="urn:ietf:params:xml:ns:netconf:base:1.0">
  <data>
      <!-- XML content/response... -->
  </data>
</rpc-reply>
This <rpc-reply> example assumes that the XML namespace is in the <rpc> sent by the client. Note that the actual data response coming from the NETCONF server is embedded within the <data> tag.

Next, we’ll show how the NETCONF request dictates which particular NETCONF operation (RPC) it’s requesting of the server.

Operations
The outermost XML element is always the type of message being sent (e.g., <rpc> or <rpc-reply>). When you are sending a NETCONF request from the client to the server, the next element, or the child of the message type, is the requested NETCONF (RPC) operation. You saw a list of NETCONF operations in Table 10-3, and now we’ll take a look at some of them.

The two primary operations we review in this chapter are <get> and <edit-config>. The <get> operation retrieves running configuration and device state information:

<rpc message-id="101" xmlns="urn:ietf:params:xml:ns:netconf:base:1.0">
  <get>
    <!-- XML content/response... -->
  </get>
</rpc>
Since <get> is the child element within the <rpc> message, this means the client is requesting a NETCONF <get> operation.

Within the <get> hierarchy, optional filter types allow you to selectively retrieve a portion of the running configuration—namely, subtree and XPath filters. Our initial focus is on subtree filters, which allow you to provide an XML document, which is a subtree of the complete XML tree hierarchy that you wish to retrieve in a given request. Later, in “Using ncclient with Cisco IOS XE”, we will use the XPath filter to inject the filter into the query path.

Example 10-6 references a specific XML data object by using the <native> element and the http://cisco.com/ns/yang/cisco-ios-xe-native URL. This data object is the XML representation of a specific data model that exists on the target device. This data model represents a full running configuration as XML, but in the example, we are requesting only the <interface> configuration hierarchy.

NOTE
As shown throughout this chapter, the actual JSON and XML objects sent can be based on either standard or vendor-specific models (often under the “native” tag).

The next two examples, for the <get> operation, are XML requests from a Cisco IOS XE device.

Example 10-6. NETCONF GET interfaces config
<rpc message-id="101" xmlns="urn:ietf:params:xml:ns:netconf:base:1.0">
  <get>
    <filter type="subtree">
      <native xmlns="http://cisco.com/ns/yang/Cisco-IOS-XE-native"> 1
        <interface></interface>
      </native>
    </filter>
  </get>
</rpc>
1
The custom xmlns, under the Cisco custom domain, represents a vendor-specific data model.

WARNING
In this edition of the book, we’ve updated the vendor model reference because the custom vendor data model changed. This may happen again in the future, so it’s necessary to always be aware of the supported capabilities.

You could add more elements to the filter’s XML tree to narrow the response that comes back from the NETCONF server. You will now add two elements to the filter—​so instead of receiving the configuration objects for all interfaces, you’ll receive the configuration of only GigabitEthernet1:

<filter type="subtree">
  <native xmlns="http://cisco.com/ns/yang/Cisco-IOS-XE-native">
    <interface>
      <GigabitEthernet>
        <name>1</name>
      </GigabitEthernet>
    </interface>
  </native>
</filter>
The next most common NETCONF operation is <edit-config>. This operation is used to make a configuration change. Specifically, this operation loads a configuration into the specified configuration data store: running, startup, or candidate. In Example 10-7, you add a static route to the running configuration.

Example 10-7. NETCONF edit-config to add a static route
<rpc message-id="101" xmlns="urn:ietf:params:xml:ns:netconf:base:1.0">
  <edit-config> 1
    <target>
      <running/>
    </target>
    <config>    2
      <configuration>
        <routing-options>
          <static>
            <route>
              <name>0.0.0.0/0</name>
              <next-hop>10.1.0.1</next-hop>
            </route>
          </static>
        </routing-options>
      </configuration>
    </config>
  </edit-config>
</rpc>
1
The <edit-config> operation is used, setting the target configuration data store with the <target> tag. If not specified, it’ll default to the running configuration.

2
Within the <config> element, we define the data-model hierarchy that we want to load onto the target data store. This structure is based on the NETCONF capabilities that are supported on a given platform.

NOTE
Vendors can implement platform-specific options. Juniper Junos OS, for example, offers options for <edit-config>. Example 10-7 uses <config>, and requires XML configuration objects for adding a static route to a Junos OS device. Junos OS also supports <config-text> within <edit-config>, which allows you to include configuration elements using text format (curly brace or set syntax).

The <edit-config> operation also supports an attribute called operation that provides more flexibility in the way a device applies the configuration object. When the operation attribute is used, it can be set to one of five values: merge, replace, create, delete, or remove. The default value is merge. If you wanted to delete the route from the previous example, you could use delete or remove; the difference is that an error occurs if you use delete when the object doesn’t exist. You could optionally use create, but an error is raised if the object already exists. Often merge is used for making configuration changes for this reason.

Finally, you could use the replace operation if you wanted to replace a given XML hierarchy in the configuration data object. In the static route example, you would use replace if you wanted to end up with just the default static route on the device; it would automatically remove all other configured static routes.

NOTE
If the operation options still seem a little confusing, don’t worry. Once we start exploring and automating devices using NETCONF later in this section, you’ll see even more examples that use various XML objects across device types for operations such as the NETCONF merge and replace.

We’ve shown what XML documents look like when using <get> and <edit-config>. The following list describes the other base NETCONF operations:

<get-config>
Retrieves all or part of a specified configuration (e.g., running, candidate, or startup).

<copy-config>
Creates or replaces a configuration data store with the contents of another configuration data store. Using this operation requires the use of a full configuration.

<delete-config>
Deletes a configuration data store (note that the running configuration data store can’t be deleted).

<lock>
Locks the configuration data-store system of a device being updated to ensure that no other systems (NETCONF clients) can make a change at the same time.

<unlock>
Unlocks a previously issued lock on a configuration data store.

<close-session>
Requests a graceful termination of a NETCONF session.

<kill-session>
Forcefully and immediately terminates a NETCONF session.

This is not an exhaustive list of NETCONF operations, but rather the core operations that each device must support in a NETCONF implementation. NETCONF servers can also support extended operations such as <commit> and <validate>. To support extended operations like these, the device must support required dependencies called NETCONF capabilities.

The <commit> operation commits the candidate configuration as the device’s new running configuration. To support the <commit> operation, the device must support the candidate capability.

The <validate> operation validates the contents of the specified configuration (running, candidate, startup). Validation consists of checking a configuration for both syntax and semantics before applying the configuration to the device.

TIP
We’ve mentioned NETCONF capabilities twice already, so let’s give a bit more context. As you know now, NETCONF supports a base set of NETCONF RPC operations. These are defined by the device as NETCONF capabilities. NETCONF capabilities are exchanged between the client and the server during connection setup, and the capabilities supported are denoted by a URL/URI. For example, every device that supports NETCONF should support the base operation from the namespace: urn:ietf:params:xml:ns:netconf:base:1.0. Additional capabilities use the form urn:ietf:params:netconf:capability:{name}:1.x, where name is the name of the capability and the way it is usually identified (without the full namespace). When we start exploring the use of NETCONF from a hands-on perspective (in Example 10-8), you’ll get to see all capabilities a given device supports.

Content
The last layer of the NETCONF protocol stack to understand is the content. This is the actual XML document that gets embedded within the RPC operation tag elements. We already showed examples of what the content could be for particular NETCONF operations.

In Example 10-6, you looked at the content that selectively requested configuration elements for the interfaces on a Cisco IOS XE device:

<native xmlns="http://cisco.com/ns/yang/cisco-ios-xe-native">
  <interface>
  </interface>
</native>
The most important point to understand about content is that it is the XML representation of a particular schema, or data model, that the device supports. We introduced schemas and data models in Chapter 8.

After this short introduction to the basic NETCONF concepts, you’re ready to start exploring what a real NETCONF API interaction looks like.

Exploring NETCONF
As you learn new APIs, it’s advantageous to learn about associated tooling that allows you to learn the API without writing any code. You saw this with cURL when learning how to use HTTP-based APIs. For NETCONF, we are going to cover how to use an SSH client that creates an interactive NETCONF session. You’ll learn how to construct a proper NETCONF request while also seeing how the device responds to a given request, without writing any code.

WARNING
Using an interactive NETCONF over SSH session as we do here is useful for learning and exploring the use of NETCONF, but it is also unintuitive, unfriendly, and fragile. For any use case outside of learning and experimentation (especially for those dealing with production infrastructure), you should instead use higher-level libraries and tools that we’ll explore later. These higher-level resources manage the minutiae of NETCONF operations and syntax much more effectively and predictably.

In the next two examples, we show you how NETCONF works on two platforms: Junos and Cisco IOS XE.

NETCONF with Junos
Let’s start with Juniper vMX running Junos as an example. You connect to the device via SSH on port 830, the default port number for NETCONF. To connect to the device, you’ll use a standard Linux ssh command:

$ ssh -p 830 ntc@vmx1 -s netconf
NOTE
Based on the vendor implementation, you may need to supply -s netconf as you SSH to the device. The -s denotes the SSH subsystem being used.

In Example 10-8, as soon as you connect and authenticate, the NETCONF server (the router) responds with a hello message that includes all of its supported NETCONF operations, capabilities, models/schemas, and a session ID.

Example 10-8. NETCONF hello response from server
<nc:hello xmlns:nc="urn:ietf:params:xml:ns:netconf:base:1.0">
  <nc:capabilities>
    <nc:capability>urn:ietf:params:netconf:base:1.0</nc:capability>
    <nc:capability>urn:ietf:params:netconf:capability:candidate:1.0</nc:capability>
    <nc:capability>urn:ietf:params:netconf:capability:confirmed-commit:1.0</nc:capability>
    <--- output omitted for brevity --->
    <nc:capability>http://xml.juniper.net/netconf/junos/1.0</nc:capability>
    <nc:capability>http://xml.juniper.net/dmi/system/1.0</nc:capability>
  </nc:capabilities>
  <nc:session-id>77470</nc:session-id>
</nc:hello>
]]>]]>
All these capabilities announce what you could do with the device via the NETCONF interface.

Once you receive the server’s capabilities, the NETCONF connection setup process starts. The next step is to send our (client) capabilities. A capabilities exchange is required to be able to send any NETCONF requests to the server.

The hello object you’re going to send to the device to complete the capabilities exchange is the following:

<?xml version="1.0" encoding="UTF-8"?>
<hello xmlns="urn:ietf:params:xml:ns:netconf:base:1.0">
  <capabilities>
    <capability>urn:ietf:params:netconf:base:1.0</capability>
    <capability>http://xml.juniper.net/netconf/junos/1.0</capability>
  </capabilities>
</hello>
]]>]]>
Note the last six characters in the preceding XML documents: ]]>]]>. These characters denote that the request is complete and can be processed. NETCONF supports two types of message separators, depending on the supported capabilities:

urn:ietf:params:netconf:base:1.0
Structured as ]]>]]> plus a newline. We use this separator in all the NETCONF examples with Junos, such as in Example 10-8.

urn:ietf:params:netconf:base:1.1
This chunked framing uses <number> and #. We use this separator in the examples in “NETCONF with Cisco IOS XE”.

However, the hello message always uses the ]]>]]>, for backward compatibility.

NOTE
In an interactive NETCONF session, you need to explicitly use the separators; when you use a library, this low-level action is abstracted (and implicit).

As you start working with the SSH client, you’ll realize it’s not like a familiar interactive CLI, although it is an interactive session. No help menu or question mark help is available. There is no man page. It’s common to think something is broken or the terminal is frozen. It’s not. If you don’t get any errors after you copy and paste XML documents into the session terminal, things are likely going well. To break out of the interactive session, you need to press Ctrl-C on your keyboard—​there is no way to safely exit the interactive NETCONF session.

Once the client responds with its capabilities, you’re ready to start sending NETCONF requests. You can use a text editor to preconstruct your XML documents.

At this point, we’ve successfully connected to the device and exchanged capabilities, and we can now issue an actual NETCONF request. Our first example will query the device for the fxp0 interface configuration. In Example 10-9, we construct the XML document in a text editor and then copy and paste it into the interactive session.

Example 10-9. NETCONF GET operation in Junos
<?xml version="1.0" encoding="UTF-8"?>
<rpc message-id="101" xmlns="urn:ietf:params:xml:ns:netconf:base:1.0">
  <get>
    <filter type="subtree">
      <configuration>
        <interfaces>
          <interface>
            <name>fxp0</name>
          </interface>
        </interfaces>
      </configuration>
    </filter>
  </get>
</rpc>
]]>]]>
As soon as you hit Enter, the request is sent to the device. As presented in Example 10-10, you’d see the XML RPC reply from the device in near real time.

Example 10-10. NETCONF GET reply
<nc:rpc-reply xmlns:nc="urn:ietf:params:xml:ns:netconf:base:1.0"
  xmlns:junos="http://xml.juniper.net/junos/18.2R1/junos"
  message-id="101" xmlns="urn:ietf:params:xml:ns:netconf:base:1.0">
  <nc:data>
    <configuration xmlns="http://yang.juniper.net/junos/conf/root"
      junos:commit-seconds="1653021086"
      junos:commit-localtime="2022-05-20 04:31:26 UTC" junos:commit-user="ntc">
      <interfaces xmlns="http://yang.juniper.net/junos/conf/interfaces">
        <interface>
          <name>fxp0</name>
          <unit>
            <name>0</name>
            <description>MANAGEMENT_INTERFACE__DO_NOT_CHANGE</description>
            <family>
              <inet>
                <address>
                  <name>10.0.0.15/24</name>
                </address>
              </inet>
            </family>
          </unit>
        </interface>
      </interfaces>
    </configuration>
    <database-status-information></database-status-information>
  </nc:data>
</nc:rpc-reply>
]]>]]>
At this point, you may have successfully performed your first request to a network device via NETCONF and received a response. The point here isn’t to do anything with it, just as you didn’t do anything with data returned with cURL. The value is that you’ve tested and validated an XML request to retrieve the configuration for interface fxp0, and now know what the response looks like, to ease you into automating devices with Python or Go.

You’ve seen one example using NETCONF <get> operations to the device. Let’s take a look at one example introducing how to use the <edit-config> operation, which is used to make a configuration change.

To see the proper way to construct an XML request for a configuration change, you are going to first issue a get request, since that will show you the structure of the complete object that needs to get sent back to the device. This is similar to knowing different CLI commands in an operating system.

NOTE
This example uses a vendor-specific data model for the interfaces instead of the standard urn:ietf:params:xml:ns:yang:ietf-interfaces because the standard one is not available. It is not present in the output of Example 10-8.

In Example 10-10, you can notice how an interface (fxp0) is configured. If you remove the inner filter part, you could get the rest of the interfaces’ configurations. You can use the data structure used for the fxp0 interface as the foundation to modify the configurations on other interfaces, to simplify the process.

Let’s make our first NETCONF configuration change by configuring the IP address of 192.0.2.1/24 on interface ge-0/0/0. To construct the object, you’ll extract the required data from your get request in Example 10-10. The two items you need to update are as follows:

Your returned object in the <data> tag will get enclosed in a <config> tag when you want to make a configuration change using the NETCONF <edit-config> operation.

The constructed object needs to specify a target data store (i.e., running, startup, or candidate) based on what the target node supports.

These changes result in the XML message in Example 10-11.

Example 10-11. NETCONF edit-config to change IP and interface description
<?xml version="1.0" encoding="UTF-8"?>
<rpc message-id="101" xmlns="urn:ietf:params:xml:ns:netconf:base:1.0">
  <edit-config>
    <target>
      <running/>
    </target>
    <config>
      <configuration xmlns="http://yang.juniper.net/junos/conf/root">
        <interfaces xmlns="http://yang.juniper.net/junos/conf/interfaces">
          <interface>
            <name>ge-0/0/0</name>
            <unit>
              <name>0</name>
              <description>Interface with changed IP</description>
              <family>
                <inet>
                  <address>
                    <name>192.0.2.1/24</name>
                  </address>
                </inet>
              </family>
            </unit>
          </interface>
        </interfaces>
      </configuration>
    </config>
  </edit-config>
</rpc>
]]>]]>
Once this XML document is built in a text editor, it can easily be copied and pasted into an active NETCONF session.

All the examples have been successful NETCONF operations so far, but the reality is that it could take a bit to get used to the XML syntax and the data models used in each platform. Luckily, NETCONF provides a useful error message to help you fix the problem. For instance, in Example 10-11, you have targeted the candidate configuration data store to change. However, if you use another data store, such as running, the NETCONF operation will complain because it’s not expected to change it directly on this platform. It requires the use of the commit operation.

NETCONF with Cisco IOS XE
To complement the Junos NETCONF example, it’s interesting to see another NETCONF implementation—in this case, for Cisco IOS XE.

You establish the SSH NETCONF session, and receive the hello from the device, as in Example 10-8:

$ ssh -p 830 ntc@csr1 -s netconf
<?xml version="1.0" encoding="UTF-8"?>
<hello xmlns="urn:ietf:params:xml:ns:netconf:base:1.0">
<capabilities>
<capability>urn:ietf:params:netconf:base:1.0</capability>
<capability>urn:ietf:params:netconf:base:1.1</capability>
<capability>http://tail-f.com/ns/netconf/actions/1.0</capability>
<capability>http://cisco.com/ns/cisco-xe-ietf-ip-deviation?...
  revision=2016-08-10</capability>
<capability>http://openconfig.net/yang/policy-types?...
  revision=2016-05-12</capability>
<capability>urn:ietf:params:xml:ns:yang:smiv2:RFC-1212?module=RFC-1212</capability>
<--- output omitted for brevity --->
All these capabilities announce what you could do with each device via the NETCONF interface. The first big difference, among platforms, is in the number of supported capabilities—there are many more in this case. This doesn’t mean that the first is less capable than the second. It means only that both support different data models and features.

It’s interesting to notice the different organizations being referenced in each output. As expected, both support the base IETF definitions for NETCONF operations, but each vendor has its own extensions. The Juniper one uses its own data models (xml.juniper.net), as does Cisco (cisco.com and tail-f.com). In this output, you can see that it also supports OpenConfig models (openconfig.net). And last, but not least, you can observe a lot of translated data models, from SMIv2 to YANG (defined in RFC 6643), to make available the same structure data used in SNMP, via NETCONF. We talk about the data models’ definitions in “Comparing NETCONF, RESTCONF, and gNMI”.

In this case, you will send a hello message announcing your intention to use base:1.1 (which you did not use in the previous NETCONF examples with Junos) because it allows the chunked framing delineator that is the only one supported in this platform:

<hello xmlns="urn:ietf:params:xml:ns:netconf:base:1.0">
    <capabilities>
        <capability>urn:ietf:params:netconf:base:1.0</capability>
        <capability>urn:ietf:params:netconf:base:1.1</capability>
    </capabilities>
</hello>
]]>]]>
To run a get-config operation, you target the source data store and use the chunked framing delineator before and after the message (the number is arbitrary):

#200
<rpc xmlns="urn:ietf:params:xml:ns:netconf:base:1.0" message-id="101">
  <get-config>
    <source>
      <running/>
    </source>
  </get-config>
</rpc>
##
This code outputs the full running configuration in XML format. You could apply filters to narrow the desired data.

TIP
We’ve stated this a few times already, but we’re going to restate it because it’s extremely important. As you get started using APIs, you need to know how to construct the proper request object. This is often challenging as you get started, but the hope is you can find easy ways to help figure out how to build these objects. This help could come from API documentation, tooling built to interface with the underlying schema definitions files such as XSDs or YANG modules, or even CLI commands on the device. For example, Cisco Nexus and Juniper Junos OS have CLI commands that show you exactly what the XML document needs to be for a given request.

As we wrap up the NETCONF section to move to other network management interfaces, keep in mind that NETCONF uses only XML encoding and SSH transport (you’re going to find different options next), and it supports data-model-based operations like the other interfaces.

After this deep dive into NETCONF, let’s continue with its cousin, RESTCONF.

Using RESTCONF
In “Getting Familiar with HTTP-Based APIs”, we explained how popular REST APIs work, and afterward, in “Using NETCONF”, we presented the benefits of exposing network management operations based on data models. The network automation community’s demand to bring both together should not come as a surprise. The answer was RESTCONF, defined in RFC 8040, which has been adopted by many NOSs.

Adhering to REST principles implies some limitations and simplifications compared to NETCONF:

RESTCONF doesn’t support network-wide transactions because it requires stateful communications. So, the clients should manage failure scenarios in the system-by-system interactions.

RESTCONF drops the data-store concept; only a single data store exists, equivalent to the running one.

No locking operation is supported.

WARNING
Don’t consider NETCONF and RESTCONF as mutually exclusive; the truth is almost the opposite. Devices can provide both interfaces relying on the same backend functionalities, so it’s up to the client to choose the most convenient one for particular use cases. However, dual interaction could lead to compatibility issues. For instance, locking a data store via NETCONF and accessing it via RESTCONF will raise an error.

As it’s a combination of REST APIs and NETCONF concepts, RESTCONF assembles characteristics from both. RESTCONF uses HTTP with the standard requests: GET, PUT, PATCH, POST, and DELETE, corresponding to some NETCONF operations, as you can see in Table 10-4.

Table 10-4. RESTCONF-to-NETCONF operation mapping
RESTCONF	NETCONF
GET

<get>/<get-config>

POST

<edit-config>> (nc:operation="create")

PUT

<edit-config>> (nc:operation="create/replace")

PATCH

<edit-config>> (nc:operation, depends on the content)

DELETE

<edit-config>> (nc:operation="delete")

RESTCONF supports JSON and XML encoding. As the network developer, you have the choice to work with whichever data format you prefer. The structure of the content data is defined using YANG models, like NETCONF.

Exploring RESTCONF in Cisco IOS XE
As in “Exploring HTTP-based APIs with cURL”, you’ll continue using cURL to explore RESTCONF because it’s just another REST API. In this section, we use the Cisco IOS XE platform to demonstrate basic RESTCONF interaction. Similar to NETCONF, you may need to activate the RESTCONF interface in your network device, if it’s supported.

TIP
A great resource to extend your knowledge about RESTCONF is the Cisco DevNet learning tracks. You will also find content about the other interfaces covered in this chapter.

Before getting started, you can use a well-known path (/.well-known/host-meta) to discover the RESTCONF path within the API:

$ curl https://csr1/.well-known/host-meta -k -u 'ntc:ntc123' 1

<XRD xmlns='http://docs.oasis-open.org/ns/xri/xrd-1.0'>      2
    <Link rel='restconf' href='/restconf'/>
</XRD>
1
The -u (or user) cURL argument provides the user and password for Basic authentication. In Basic authentication, the user and password are sent in clear text over the network (Base64 encoded), so it’s not a recommended method for production.

2
By default, the API returns the response in XML. Remember, you could change this behavior by using the HTTP Accept header.

WARNING
A Base64-encoded string does not mean it has been encrypted. You can easily encode and decode Base64-encoded strings with Python using the base64 BY Python module:

>>> import base64
>>>
>>> encoded = base64.b64encode('ntc:ntc123')
>>> encoded
'bnRjOm50YzEyMw=='
>>>
>>> text = base64.b64decode(encoded)
>>> text
'ntc:ntc123'
>>>
Now, knowing that the RESTCONF API is located at /restconf, you are ready to start exploring it:

$ curl -k -X GET https://csr1/restconf \
    -H 'Accept: application/yang-data+json' \ 1
    -u 'ntc:ntc123'

{
  "ietf-restconf:restconf": {
      "data":{},                              2
      "operations":{},                        3
      "yang-library-version":"2016-06-21"
  }
}
1
Using the Accept header with application/yang-data+json, the response comes in JSON format instead of the default XML.

2
The data path will contain all data resources.

3
The operations path will include the data-model-specific operations.

Then, similar to NETCONF, you can discover the capabilities offered via the RESTCONF interface. In Example 10-8, you can target the path /restconf/data/netconf-state/capabilities to obtain the very same list you got in “NETCONF with Cisco IOS XE”:

$ curl -k -X GET https://csr1/restconf/data/netconf-state/capabilities \
    -H 'Accept: application/yang-data+json' \
    -u 'ntc:ntc123'

{
  "ietf-netconf-monitoring:capabilities": {
    "capability": [
      "urn:ietf:params:netconf:base:1.0",
      "urn:ietf:params:netconf:base:1.1",
      "http://cisco.com/ns/yang/Cisco-IOS-XE-native?module=Cisco-IOS-XE-native&
        revision=2019-11-01", 1
      # output omitted for brevity
    ]
  }
}
1
This model contains all the Cisco IOS XE native configuration models.

In Example 10-12, we explore Cisco-IOS-XE-native:native under the /restconf/data/ path.

Example 10-12. GET configuration with RESTCONF
$ curl -k -X GET https://csr1/restconf/data/Cisco-IOS-XE-native:native \
    -H 'Accept: application/yang-data+json' \
    -u 'ntc:ntc123'

{
  "Cisco-IOS-XE-native:native": {
    "version": "17.1",
    "memory": {
      "free": {
        "low-watermark": {
          "processor": 72107
        }
      }
    },
    # output omitted for brevity
  }
}
The next subsections explore the RESTCONF API in more detail, including updating the configuration (with PATCH and PUT operations) and executing operations.

Updating configuration via RESTCONF
As in any other REST API, you can change the content of the API via POST, PATCH, PUT, and DELETE operations. In this case, you want to add two OSPF network statements to an existing OSPF configuration, so the PATCH operation will allow updating it starting from the current state.

As we’ve already said, when dealing with APIs, you need to know the expected data model structure that the API can understand. You can check the data model specification from the announced capabilities, but it’s quicker to explore the data model with a GET operation (as in Example 10-12). In the previous output, you got the full configuration, but to check a specific section, you can append it to the previous path. For instance, by adding router/ to the previous path, you get only the router section that contains the OSPF current configuration:

"router": {
  "Cisco-IOS-XE-ospf:router-ospf": {
    "ospf": {
      "process-id": [
        {
          "id": 10,
          "network": [
            {
              "ip": "192.0.2.0",
              "wildcard": "0.0.0.7",
              "area": 0
            },
            {
              "ip": "192.0.2.64",
              "wildcard": "0.0.0.7",
              "area": 0
            }
          ],
          "router-id": "192.0.2.1"
        }
      ]
    }
  }
}
Now, in Example 10-13, we request a PATCH operation with the data payload containing two new networks to update the OSPF configuration.

Example 10-13. RESTCONF PATCH to add OSFP statements
$ curl -k
    -X PATCH "https://csr1/restconf/data/Cisco-IOS-XE-native:native/router/router-ospf" \ 1
    -H 'Content-Type: application/yang-data+json' \
    -H 'Accept: application/yang-data+json' \
    -u 'ntc:ntc123' \
    -d $'{
  "router-ospf": {
    "ospf": {
      "process-id": [
        {
          "id": 10,                                                                       2
          "network": [
            {
              "ip": "192.0.2.128",
              "wildcard": "0.0.0.7",
              "area": 0
            },
            {
              "ip": "192.0.2.192",
              "wildcard": "0.0.0.7",
              "area": 0
            }
          ]
        }
      ]
    }
  }
}'
1
We target the router-ospf leaf in the URL path.

2
The process-id is the same as the previous one (10), so it will append the new networks to the existing ones. We can check it by repeating the previous GET operation, and it would show the four OSPF networks.

If instead of PATCH you use the PUT method, the whole router-ospf will be replaced by the new one. This method allows much more efficient configuration management than traditional CLI configuration. In the CLI, adding configurations is trivial. However, removing or negating commands is complex. For example, what if you have a single instance of OSPF running with 50 network statements, but because of a change in design, you need only 2 network statements? You will have to know which 48 statements need to be negated with a no command. This process is arduous and mundane as you extrapolate the effort for all types of configurations on a network device. For our example, wouldn’t it be easier to take the opposite approach—focus on the configuration that should exist on the network device? This is a growing trend, becoming more possible thanks to newer APIs and ways of thinking. This is called declarative configuration.

NOTE
In “Updating configuration via RESTCONF with net/http”, you can see an example of using the HTTP PUT method to update the configuration.

Understanding the YANG PATCH HTTP operation
Following REST principles, RESTCONF comes with two limitations when compared to NETCONF. First, HTTP calls should not carry the state in between, so a transaction is limited to one HTTP call. Second, an HTTP request implements one type of create, read, update, and delete (CRUD) operation specified by the HTTP method (e.g., GET or POST). The second limitation can be overcome by a new YANG PATCH HTTP media type that allows combining various operation types on the same HTTP request.

Before you start using it, you need to validate that the feature (i.e., urn:ietf:params:restconf:capability:yang-patch:1.0) is supported in the target platform. You can look for this capability in the restconf-state/capabilities endpoint:

$ curl -k -X GET \
    https://csr1/restconf/data/ietf-restconf-monitoring:restconf-state/capabilities \
    -H 'Accept: application/yang-data+json' \
    -u 'ntc:ntc123'

{
  "ietf-restconf-monitoring:capabilities": {
    "capability": [
      # output omitted for brevity
      "urn:ietf:params:restconf:capability:yang-patch:1.0", 1
    ]
  }
}
1
Notice the yang-patch capability listed under the restconf-state/capabilities endpoint.

To illustrate how YANG PATCH works, you will manage loopback interfaces, adding and removing them. With YANG PATCH, you will first create a Loopback0, and in a second request, you will remove it and add a new Loopback1.

In Example 10-14, you have the YANG PATCH payload necessary to create a new Loopback0 interface. Remember, you need to know the expected data structure for the value before.

Example 10-14. YANG PATCH with one operation
<yang-patch xmlns="urn:ietf:params:xml:ns:yang:ietf-yang-patch">
  <patch-id>add-Loopback0-patch</patch-id> 1
  <edit>                                   2
    <edit-id>edit1</edit-id>
    <operation>create</operation>          3
    <target>/Loopback=0</target>           4
    <value>                                5
        <Loopback xmlns="http://cisco.com/ns/yang/Cisco-IOS-XE-native">
          <name>0</name>
        </Loopback>
    </value>
  </edit>
</yang-patch>
1
Unique ID to identify the YANG PATCH request.

2
A YANG PATCH is an ordered list of edits, each one with an identifier (edit-id).

3
The type of operation: create, delete, insert, merge, move, replace, or remove.

4
Specifies the target node the operation targets.

5
This is optional (for example, a move operation doesn’t need it), but for a create operation, it is the content to be created in the target.

Then, you can use this data payload with an HTTP PATCH operation with a special Content-Type, application/yang-patch+xml. The @ specifies using a file instead of a data payload directly:

$ curl -k -X PATCH "https://csr1/restconf/data/Cisco-IOS-XE-native:native/interface" \
    -H 'Content-Type: application/yang-patch+xml' \
    -H 'Accept: application/yang-data+xml' \
    -u 'ntc:ntc123' \
    -d '@create-loopback-0.xml'

<yang-patch-status xmlns="urn:ietf:params:xml:ns:yang:ietf-yang-patch">
  <patch-id>add-Loopback0-patch</patch-id>
  <ok/>
</yang-patch-status>
NOTE
Remember that YANG PATCH is not an HTTP method but a media type.

Example 10-14 helps you understand the key parts of the YANG PATCH media type, but a single operation can be done without this special type. In Example 10-15, we add a new interface, Loopback1, and remove the previous one, Loopback0, in the same HTTP request.

Example 10-15. YANG PATCH with two operations
<yang-patch xmlns="urn:ietf:params:xml:ns:yang:ietf-yang-patch">
  <patch-id>add-remove-loopback-patch</patch-id>
  <edit>
    <edit-id>edit1</edit-id>
    <operation>create</operation>
    <target>/Loopback=1</target>
    <value>
        <Loopback xmlns="http://cisco.com/ns/yang/Cisco-IOS-XE-native">
          <name>1</name>
        </Loopback>
    </value>
  </edit>
  <edit>
    <edit-id>edit2</edit-id>
    <operation>remove</operation>
    <target>/Loopback=0</target>
  </edit>
</yang-patch>
TIP
The order of the operations within a YANG PATCH payload is not relevant. What is important is the final configuration outcome state and its behavior as an atomic transaction, for one single device. Either all the operations work well or the whole set of changes is rejected.

The YANG PATCH media supersedes all HTTP methods. This means that you could concentrate all your CRUD operations via YANG PATCH if you wish.

Discovering RESTCONF operations
RESTCONF, like NETCONF, supports managing modeled data but also running operations. You can list the supported operations in the /restconf/operations path:

$ curl -k -X GET https://csr1/restconf/operations \
    -u 'ntc:ntc123'  \
    -H 'Accept: application/yang-data+json' \
    | python3 -m json.tool

{
    "ietf-restconf:operations": {
        "Cisco-IOS-XE-rpc:factory-reset":                  1
          "/restconf/operations/Cisco-IOS-XE-rpc:factory-reset",
        "ietf-event-notifications:establish-subscription": 2
          "/restconf/operations/ietf-event-notifications:establish-subscription",
        "ietf-event-notifications:create-subscription":    2
          "/restconf/operations/ietf-event-notifications:create-subscription",
        # output omitted for brevity
    }
}
1
Notice the various sources. Some are vendor specific, defined by Cisco in this case.

2
Other operations are defined by the IETF.

With this, we conclude the exploration of the RESTCONF interface. In “The Go net/http Package”, we will come back to RESTCONF, but will use the Go net/http package to automate it. After NETCONF and RESTCONF, it’s time to explore another network management interface: gNMI.

Using gRPC and gNMI
The gRPC Network Management Interface (gNMI) is, as its name indicates, a network management interface built on top of gRPC. It tries to overcome some limitations seen in the network management space, from SNMP to NETCONF, solving the two network management goals—configuration management and state retrieval—in a data-model-oriented way. Like NETCONF or RESTCONF, gNMI is a data model interface. However, it has characteristics that make it different:

It uses gRPC as the transport protocol, instead of SSH or HTTP, and protobuf for encoding.

It is defined and maintained by the OpenConfig consortium, led by Google. In contrast, NETCONF and RESTCONF are IETF standards.

gNMI is an alternative to other network management protocols that has gained popularity because of a fast initial development, and a capable and simple feature set. Most popular NOSs support it: Nokia Service Router Operating System (SROS) and SR Linux; Cisco IOS XR, IOS XE, and NX-OS; Arista EOS; Junos OS; and SONiC. We compare gNMI and the other data model interfaces in more detail in “Comparing NETCONF, RESTCONF, and gNMI”.

Complementary to gNMI, the gRPC Network Operations Interface (gNOI), defines operational commands on network devices (e.g., ping or reboot). We don’t cover gNOI in this book, but it follows the same principles as gNMI.

But before getting into the gNMI capabilities, we need to introduce the framework that makes it possible: gRPC.

Understanding gRPC
gRPC (initially called Stubby) was created by Google and evolved into a public project within the CNCF. Before getting into the characteristics of gRPC, let’s understand the motivations behind its creation.

Google’s application stack is built around distributed microservices. Each application can be implemented in a different programming language and requires high performing communications to execute RPC operations. gRPC was designed to support these requirements by being useful in multiple environments, including intra-data-center and end-user applications. gRPC offers low latency and extensibility to add extra features, such as load balancing or tracing.

gRPC design embraces several principles, and the following are a few of the most significant ones:

High performance
Fast and efficient communication between services is one core principle of gRPC design. For example, one significant change versus REST APIs is the use of static paths instead of dynamic ones. Dynamic paths include query parameters that need to be parsed before processing the call, which makes it slower and more complex. In gRPC, everything is part of the message body.

Payload agnostic
The framework supports multiple content types for data serialization and encoding, such as Protocol Buffers, JSON, and XML. Because of the high-performance requirement, the most common one is Protocol Buffers. In this section, we will use this data format, introduced in Chapter 8.

Several communication patterns
gRPC allows various communication patterns, from traditional request/response (unary) to unidirectional and bidirectional streaming. Also, it can work in asynchronous or synchronous mode to enable scalability and streaming processing.

Language independent
gRPC clients and servers can be built in multiple popular languages, such as Python and Go, and also in cross-platform environments. In “A gRPC example”, we demonstrate this by implementing the client in Python, and the server in Go.

You can extend your gRPC knowledge at https://grpc.io.

Using gRPC as a transport protocol brings a lot of flexibility. It works over TCP without a predefined port—the port is defined by each application. Also, it comes without any predefined calls and messages. It’s up to each application to define them. You will see a specific implementation for network management in “Understanding the gNMI interface”.

The best way to understand a protocol like gRPC is via an example, so let’s dig in.

![Alt text](image-2.png)\
\
Figure 10-3 shows a gRPC service example using a Python client and a Go server. The gRPC service is defined by the protobuf file from Example 8-12. All the files for these examples are located at https://github.com/oreilly-npa-book/examples/tree/v2/ch10-apis/grpc.

npa2 1003
Figure 10-3. gRPC communication
The protobuf file (networkstuff.proto, defined in Chapter 8) contains the definition of the messages (the data types) and the service (the RPC operation). The only addition in comparison to Example 8-12 is the option go_package = "./networkstuff"; to enable auto-generation of code for Go.

Running a gRPC server in Go
One key feature of gRPC is its ability to create code bindings from the protobuf definition. This means that the messages and services will be transformed into type structures and functions, respectively. In the directory, the files have been autogenerated, so you don’t need to do it. The only Go file manually defined is server.go; the rest of Go files in the folder are autogenerated by the protobuf compiler.

NOTE
To autogenerate gRPC bindings, you need to install the protocol buffer compiler. The installation for various platforms is described at https://oreil.ly/1NtLG.

We define and initialize a new Go module, grpc_example, that uses the autogenerated module networkstuff. The server.go file (part of the grpc_example module) uses the generated bindings to expose the gRPC service, which we analyze here in three parts.

Example 10-16 contains the beginning of the server.go file that will perform as the gRPC server.

Example 10-16. gRPC Go server.go, part 1
package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "net"

    // implementation of gRPC for Go maintained by Google
    "google.golang.org/grpc"

    // our own package that was built before, and updated to point
    // to ./networkstuff via a replacement in the go.mod file
    pb "github.com/pkg/networkstuff"
)

// custom struct extending the pb.RouterServiceServer plus a
// local cache of Routers
type routerServiceServer struct {
    pb.UnimplementedRouterServiceServer
    localRouters []*pb.Router
}

// method necessary to match the interface definition for
// pb.RouterServiceServer, with the same signature
func (s *routerServiceServer) GetRouter(ctx context.Context,
    router_request *pb.RouterRequest) (*pb.Router, error) {
    // Use the local cache of Routers to match by the
    // router identifier
    for _, router := range s.localRouters {
        if router.Id == router_request.Id {
            return router, nil
        }
    }
    // No router was found, return a nameless router
    return &pb.Router{}, nil
}

// -> continues to part 2
The first part of the file contains the necessary imports, with special attention given to grpc and networkstuff. Then, the file defines the routerServiceServer struct, which is implementing the RouterServiceServer interface (from networkstuff_grpc.pb.go), plus a list of Routers objects to serve as a poor-man data store. The GetRouter method is required by the interface (Chapter 7 provides more details about implementing struct interfaces), and it returns the corresponding element from the localRouters list.

NOTE
The gRPC package is part of Go included packages at https://oreil.ly/7LeGx, but it can also be found at https://oreil.ly/WTNma.

Continuing with the server.go file, Example 10-17 includes a server variable with all the data that the gRPC server will serve. The data could have been loaded from an external database or a JSON file, but we decided to directly initialize the objects.

Example 10-17. gRPC Go server.go, part 2
// -> comes from part 1
// server contains the data to expose via grpc
var server = &routerServiceServer{
    localRouters: []*pb.Router{
        &pb.Router{
            Id:       1,
            Hostname: "Router A",
            Interfaces: []*pb.Interface{
                &pb.Interface{
                    Id:          1000,
                    Description: "Gi 0/0/0",
                },
                &pb.Interface{
                    Id:          1001,
                    Description: "Gi 0/0/1",
                },
            },
        },
  // omitted for brevity
  },
}
// -> continues to part 3
Finally, in Example 10-18, a new TCP socket is defined and initialized to start serving TCP connections.

Example 10-18. gRPC Go server.go, part 3
// -> comes from part 2

func main() {
    // Create a TCP server listener in 50051 port
    lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 50051))
    // This is the common Go Pattern to handle errors
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }
    // Bootstrap a gRPC server with defaults
    var opts []grpc.ServerOption
    grpcServer := grpc.NewServer(opts...)

    // Register the custom RouterServiceServer implementation to
    // the gRPC server
    pb.RegisterRouterServiceServer(grpcServer, server)
    // Attach the gRPC server to the TCP port 50051 opened before
    // to start serving requests.
    grpcServer.Serve(lis)
}
With the server.go file ready, you can start the gRPC server and move to another terminal while the server is listening for gRPC connections:

ch10-apis/grpc$ go run server.go
Running a gRPC client in Python
gRPC can generate code bindings in multiple languages, so to complete the gRPC example, you will use Python to run the client gRPC requests. First, as we did for Go, it’s necessary to autogenerate the Python bindings with the libraries grpcio and grpcio-tools. Using these libraries, we already generated the Python files that are available in the examples located at https://github.com/oreilly-npa-book/examples/tree/v2/ch10-apis/grpc.

In Example 10-19, we demonstrate in the Python interpreter how to use the generated gRPC code to establish a client request toward the gRPC server (running in the other terminal).

Example 10-19. Python gRPC client usage
>>> import grpc                                        1
>>> import networkstuff_pb2_grpc as pb2_grpc           1
>>> import networkstuff_pb2 as pb2                     1
>>> channel = grpc.insecure_channel("localhost:50051") 2
>>> stub = pb2_grpc.RouterServiceStub(channel)         3
>>> router_request = pb2.RouterRequest(id=1)           4
>>> result = stub.GetRouter(router_request)            5
>>> type(result)
<class 'networkstuff_pb2.Router'>
>>> result                                             6
id: 1
hostname: "Router A"
interfaces {
  id: 1000
  description: "Gi 0/0/0"
}
interfaces {
  id: 1001
  description: "Gi 0/0/1"
}
1
Imports the grpc library and the autogenerated bindings from the protobuf.

2
Establishes a grpc insecure channel (a TCP connection) without certificate validation.

3
Defines a stub (a single gRPC client on top of the TCP connection) by using the protobuf specification.

4
Creates a proper request by using the protobuf specs embedded in a Python class and the id of the router to check.

5
Using the stub and the router_request, we can call the predefined GetRouter method.

6
Finally, we check that the response serialized with the protobuf specs and is of type Router.

We leave to you the exercise to create other gRPC requests, with existing IDs (id=2) or with nonexistent ones, to see how the client behaves. The expectations are to get Router B, or an empty router for another ID.

After this quick overview of gRPC, we start digging into the gNMI interface, which uses gRPC for network management.

Understanding the gNMI interface
gNMI defines a small set of gRPC operations: CapabilityRequest, GetRequest, SetRequest, and SubscribeRequest. These operations are defined in the gnmi.proto file, which is part of the OpenConfig project and can be found at https://oreil.ly/pul-l.

In “Comparing NETCONF, RESTCONF, and gNMI”, we dive into the differences between gNMI versus NETCONF and RESTCONF, but all these data-model-oriented interfaces share a lot of similarities, as you can see in Table 10-5, which compares the NETCONF and gNMI basic operations. More details about available gNMI operations can be found at the gNMI specification documentation.

Table 10-5. Mapping gNMI operations to NETCONF
gNMI	NETCONF
CapabilityRequest

hello

GetRequest

get/get-config

SetRequest

edit-config

SubscribeRequest

establish-subscription

gNMI has gained a lot of popularity in the industry because of the open source contributions by the OpenConfig consortium (mostly by Google), and in the observability realm because it provided the first implementation of dial-in streaming telemetry. You’ll learn more about this in “Understanding model-driven telemetry”.

After this short introduction to gNMI, let’s get hands-on and explore it.

Exploring gNMI with gNMIc
To explore gNMI, you will use the gNMIc CLI client. This open source client has been adopted under the OpenConfig umbrella and can be found at https://oreil.ly/d343N.

NOTE
In this section, you use the gNMIc CLI version, but in “The OpenConfig gNMIc Go Package”, you will use the base Go package to create network automation scripts for interacting with gNMI.

First, you need to install gNMIc. You can use a downloaded bash script for macOS or Linux (or you could use gNMIc as a Docker container). For more information, check https://gnmic.openconfig.net/install; here’s the command:

$ bash -c "$(curl -sL https://get-gnmic.openconfig.net)"
We use an Arista switch to explore gNMI. Because there is no standard TCP port for gNMI, the first thing to double-check is the TCP port being used (you can check it in the platform documentation or directly on the platform settings). The default port for Arista is TCP/6030, but you can change it if necessary.

Then, in Example 10-20, you can check which data models are supported behind the gNMI interface, using the CapabilityRequest operation (capabilities), equivalent to the NETCONF hello.

Example 10-20. gNMI capabilities
$ gnmic -a eos-spine1:6030 -u ntc -p ntc123 --insecure capabilities
gNMI version: 0.7.0
supported models:
  - openconfig-network-instance-types, OpenConfig working group, 0.9.3
  - iana-if-type, IANA,
  - openconfig-vlan, OpenConfig working group, 3.2.1
  - arista-vlan-deviations, Arista Networks <http://arista.com/>, 1.0.2
  - ietf-yang-types, IETF NETMOD (NETCONF Data Modeling Language) Working Group,
  // omitted models for brevity
supported encodings:
  - JSON
  - JSON_IETF
  - ASCII
Example 10-20 introduces two important topics to discuss:

Supported data models
The supported data models define what can be changed/retrieved from the network device. In the output, you can notice data models coming from standards bodies (IANA and IETF), the OpenConfig consortium, and the vendor itself. Also, notice the versioning, relevant to understand the status of its definition.

Supported encoding
This is the encoding of the payload (not the transport encoding, which is protobuf). It supports JSON and JSON_IETF (to support some serialization of YANG models), which allow structured data. It also supports ASCII encoding, used for semistructured CLI configuration, depending on the gNMI backend implementation.

You can find more information about these two topics at the gNMI specification on GitHub.

gNMI GetRequest
Once you know the interface capabilities, you can do your first gNMI Get operation. You can retrieve the actual configuration state from the running configuration, targeting the config container in an OpenConfig data model: /interfaces/interface/config. gNMI paths are a simplified variant of the XPath syntax. These paths can be obtained by observing the YANG data models or by using tools to extract paths from them—e.g., Cisco YANG Suite.

NOTE
OpenConfig and IETF advocated for different styles of data model organization. In OpenConfig, a model node has two containers, config and state, to represent the intended and the operational state, respectively. In contrast, The IETF proposed to split both configurations into different data stores. There is more information in Rob Shakir’s blog post “OpenConfig and IETF YANG Models: Can they converge?”.

In Example 10-21, the configuration contains only a few leaves (mtu, name, and type), because only a few items are configured (the defaults are not shown). This is something that depends on each data-model specification.

Example 10-21. Get config interface with gNMIc
$ gnmic -a eos-spine1:6030 -u ntc -p ntc123 --insecure  --gzip \
    get \
    --path '/interfaces/interface/config'

[
  {
    "source": "eos-spine1:6030",
    "timestamp": 1664428366933949209,
    "time": "2022-09-29T05:12:46.933949209Z",
    "updates": [
      {
        "Path": "interfaces/interface[name=Management0]/config",
        "values": {
          "interfaces/interface/config": {
            "openconfig-interfaces:mtu": 0,
            "openconfig-interfaces:name": "Management0",
            "openconfig-interfaces:type": "iana-if-type:ethernetCsmacd"
          }
        }
      }
    ]
  }
]
TIP
In this example, you see only one interface, Management0. If more were present, you could get the same behavior using XPath (filtering by the interface name) and target this specific interface with the path interfaces/interface[name=Management0]/config.

Now, instead of targeting the config container, we target the state one. Besides the interface configuration (with default values), the interface stats will be returned:

$ gnmic -a eos-spine1:6030 -u ntc -p ntc123 --insecure --gzip \
    get \
    --path '/interfaces/interface/state'
[
  {
    "source": "eos-spine1:6030",
    "timestamp": 1664428346403170400,
    "time": "2022-09-29T05:12:26.4031704Z",
    "updates": [
      {
        "Path": "interfaces/interface[name=Management0]/state",
        "values": {
          "interfaces/interface/state": {
            "arista-intf-augments:inactive": false,
            "openconfig-interfaces:admin-status": "UP",
            "openconfig-interfaces:counters": {
              "in-unicast-pkts": "1694",
              "out-unicast-pkts": "4410"
              // omitted output for brevity
            },
            "openconfig-interfaces:ifindex": 999999,
            "openconfig-interfaces:last-change": "1664343413005553920",
            "openconfig-interfaces:mtu": 0,
            "openconfig-interfaces:name": "Management0",
            "openconfig-interfaces:oper-status": "UP",
            "openconfig-interfaces:type": "iana-if-type:ethernetCsmacd"
          }
        }
      }
    ]
  }
]
gNMI SetRequest
Next, use the SetRequest operation to change the configuration of the interface. As you did in NETCONF and RESTCONF, when updating data-model-based configurations, you need to understand the data model. You can check the OpenConfig interfaces data model, and you will discover that it supports a description attribute that was not retrieved in the previous Get examples because it was not set yet, and there is no default value. Our goal is to define the interface description to New Description:

$ gnmic -a eos-spine1:6030 -u ntc -p ntc123 --insecure --gzip \
    set \
    --update-path '/interfaces/interface[name=Management0]/config/description' \
    --update-value 'New Description'
Notice the use of set instead of get and the new parameters: update-path and update-value. In this case, we use XPath filtering to target a specific interface (if not, you would be updating all the interfaces with the same description) and the corresponding value. When the SetRequest is executed successfully, you receive a response confirming the operation.

Now, if you retrieve the configuration as you did in Example 10-21, the new attribute will show up:

[
  {
    "source": "eos-spine1:6030",
    "timestamp": 1664428969195249733,
    "time": "2022-09-29T05:22:49.195249733Z",
    "updates": [
      {
        "Path": "interfaces/interface[name=Management0]/config",
        "values": {
          "interfaces/interface/config": {
            "openconfig-interfaces:description": "New Description",
            "openconfig-interfaces:mtu": 0,
            "openconfig-interfaces:name": "Management0",
            "openconfig-interfaces:type": "iana-if-type:ethernetCsmacd"
          }
        }
      }
    ]
  }
]
However, not all the attributes are configurable because of dependencies on the platform. For instance, the mtu attribute is not configurable in this Arista platform, and if you try to update it, you will receive an error, as you can see next:

$ gnmic -a eos-spine1:6030 -u ntc -p ntc123 --insecure --gzip \
  set \
  --update-path '/interfaces/interface[name=Management0]/config/mtu' \
  --update-value '1400'

target "eos-spine1:6030" set request failed: target "eos-spine1:6030" SetRequest failed:
  rpc error: code = Aborted desc = failed to apply: Unavailable command (not supported
  on this hardware platform) (at token 1: 'mtu'): CLI command 3 of 5 'l2 mtu 1400'
  failed: invalid command
CLI Commands:
1 interface Management0
2 l2 mtu 1400
3 exit

Error: one or more requests failed
NOTE
This error message also shows something interesting—what is happening behind the scenes. Every platform implements the RPC operations differently. In this case, you see from the CLI commands executed that the device is translating the YANG data model changes into CLI commands.

gNMI Subscribe
gNMI was the first data-model interface to support data-model-driven telemetry (dial-in) with the subscribe operation. We go deep on telemetry in “Understanding model-driven telemetry”, but here you have an example of subscribing to the interface counters, targeting the counters’ path in the prefix:

$ gnmic -a eos-spine1:6030 -u ntc -p ntc123  --insecure \
  subscribe \
  --path "/interfaces/interface/state/counters"

{
  "source": "eos-spine1:6030",
  "subscription-name": "default-1664428777",
  "timestamp": 1664428707256105665,
  "time": "2022-09-29T05:18:27.256105665Z",
  "prefix": "interfaces/interface[name=Management0]/state/counters",
  "updates": [
    {
      "Path": "in-octets",
      "values": {
        "in-octets": 141271
      }
    },
    {
      "Path": "in-unicast-pkts",
      "values": {
        "in-unicast-pkts": 1779
      }
    }
  ]
}
# other outputs omitted for brevity
The subscribe operation keeps a session established, from the client, to get updates on the data model. In this case, we are subscribed to updates on the interface counters, but you could also subscribe to the interface operational status or the BGP session establishment, and so on. The updates can be sampled periodically or triggered when the value changes, so you get notified only when something happens.

As you may have already concluded after reviewing NETCONF, RESTCONF, and gNMI, these protocols and interfaces share common concepts but also have differences. To make this content easier to digest, we spell out the differences next.

Comparing NETCONF, RESTCONF, and gNMI
At this point, you likely have questions about the differences between the three interfaces using a data-model-driven approach (NETCONF, RESTCONF, and gNMI). In this section, we want to give you more insight into how they are related.

Everything started with NETCONF (the first RFC published in 2006) trying to address the limitations of SNMP to manage network configurations. At that moment, SNMP was widely used to monitor network’s operational state, but it was not adopted for network management. Because of this, as we explained in “Using NETCONF”, new ideas were introduced: multiple data stores, RPC operations, effective config management in transactions, and the use of data models to update/retrieve configuration data and retrieve operational data.

Notice that we stress configuration management as the primary use case for NETCONF. Implicitly, it was acknowledging that the SNMP monitoring approach, ported to NETCONF, was good enough. However, later, new concerns about a better way to retrieve operational data appeared, and the IETF started (around 2014–2015) to get requirements to implement streaming telemetry (a continuous and customized stream of data from a YANG data store), which we present in “Understanding model-driven telemetry”.

Around the same time (2014), the OpenConfig consortium was founded and led by Google, focused on implementing streaming telemetry as one of the main drivers of gNMI, a new network management interface defined as an open source project (with the first commit in 2017), instead of an internet standard. For the rest of the functionalities, it used NETCONF ideas as a reference with a simpler implementation.

More or less in parallel (2017), the IETF created the RESTCONF interface that brought together the NETCONF approach and the simpler, well-known RESTful API paradigm to promote more adoption of data-model-driven management when not all the NETCONF requirements are needed.

Table 10-6 compares the three interfaces by focusing on their main differences.

Table 10-6. Comparing data-model-driven interfaces
NETCONF	RESTCONF	gNMI
Encoding

XML

JSON or XML

protobuf or JSON

Transport

SSH

HTTP/TLS

gRPC over HTTP/2

Transaction scope

Network-wide

Single-target, single-shot

Single-target, single-shot, sequenced

Let’s take a closer look at these three main differences:

Encoding
NETCONF used the most popular encoding when it was defined (XML), and RESTCONF, even though still supporting XML (likely for reusability from NETCONF scripts), promoted JSON as a more popular encoding (and easier to read). However, gNMI chose protobuf because this binary encoding reduced the payload when compared to XML. Streaming operational data was an important use case for this interface, and this creates a lot of network traffic.

Transport
NETCONF can support multiple transport protocols, but the most common one when it was defined was SSH, and it became the de facto one. RESTCONF leveraged HTTP to become yet another REST API, and gNMI adopted the Google internal protocol gRPC to support protobuf encoding.

Transaction scope
NETCONF came from operator best practices whereby the network itself is used to test and deploy configurations that can affect multiple devices, so the network-wide transaction scope was implemented. RESTCONF, to adhere to REST principles, doesn’t support state and can deal with only one target at a time, without any relevant operation sequence. gNMI is the simplest of all, assuming a specific order of all the operations. This simple approach comes from software-oriented teams managing infrastructure, where the configuration validation is most likely done out of the box, in a development environment.

Next, so you can better understand the interfaces, we compare the development lifecycle of the organizations behind them, which are also developing vendor-neutral data models.

Network interfaces development lifecycle
NETCONF and RESTCONF are defined and promoted by the IETF. On the other side, gNMI is developed under the umbrella of the OpenConfig community. Each group has a completely different way of working, which affects how these protocols and interfaces are defined, developed, and adopted:

IETF
Its mission—defined in RFC 3936—is to develop technical recommendations to design, use, and manage the internet. This is done via an open process, with volunteers, and rough consensus from multiple parts.

OpenConfig
This working group of network operators and vendors is focused on building a vendor-independent software layer for managing network devices. It operates as an open source project, with direct contributions.

Both have the common goal to create better ways to manage networks. However, the difference is in the way they achieve it. IETF’s process of proposal, approval, and implementation usually takes longer than the same process under OpenConfig. In open source projects with strong leadership, the proposal, review, and adoption process leads to a faster release cycle. A standard addresses diverse use cases and will be stable for a while, but when you work on an open source project, you can move quickly to solve more concrete use cases and leave the door open for future changes. This capacity to deliver, especially for streaming telemetry, was one of the key factors for the adoption and popularity of the gNMI interface.

Having multiple options to solve a problem is nothing new in the networking world. Successful nonstandard solutions have been adopted multiple times, and afterward, these solutions (with small differences) were adopted by IETF (e.g., NetFlow and Internet Protocol Flow Information Export, or IPFIX). The experience says that both could likely coexist, solving different usage approaches.

OpenConfig and IETF are working on more than network management interfaces. A main focus of both organizations is to implement the right data models to describe the network information needed for configuration and operational validation. These models should work well under any of the interfaces—for instance, using an OpenConfig data model under the NETCONF interface, as you saw in “NETCONF with Cisco IOS XE”. However, as we mentioned, the structure of the data models differ in how the intended and operational data is organized.

Last, we will delve into streaming telemetry.

Understanding model-driven telemetry
Current large-scale network architectures have uncovered the operational limitations of SNMP for network monitoring. Even though SNMP is still in use, its model to retrieve data via a poll method, with the management server asking for data every time, was adding delay and extra processing to get data continuously.

These performance limitations, together with the adoption of data-model management (e.g., NETCONF or gNMI), led to the definition of model-driven telemetry. It adopts a push model to stream data from network devices continuously, providing near real-time access to operational data. Now, instead of asking for data, the management applications can subscribe to specific data from the supported data models defined with YANG. The retrieved data is structured and can be published at a defined cadence or as it changes.

With model-driven telemetry, you decide which data you need, when you need it, and where to send (or receive) it. You can subscribe to telemetry streams in two ways: dial-in and dial-out (represented in Figure 10-4).

Dial-in
This is a dynamic model: an external application opens a session to the network device and establishes one or more subscriptions (to different parts of the data store) over the same session. The network devices send operation data for as long as the session stays up.

Dial-out
This is a configured model: the subscriptions are configured on the network device in advance using any of the available interfaces, and it’s up to the device to open a telemetry session to the receiver. If this session goes down, the network device will open a new one.

![Alt text](image-3.png)

Figure 10-4. Model-driven telemetry (npa2 1004)
gNMI implemented the first streaming telemetry via the dial-in model, and it’s still the most widely adopted one. However, the dial-out model adds some benefits:

Reduces a network device’s exposure to external threats as the connection is initiated from the device itself and avoids firewall configurations to let access in.

Collectors can be stateless; they need to only listen and store the collected data. The control is on the configuration management system that created the subscription on the network device.

Nowadays, both NETCONF and gNMI support the dial-in model. Dial-out mode is simply a configuration setting, which can be done via any configuration interface (including CLI), and the data is exported over a transport protocol (TCP or UDP).

The first model-driven telemetry implementation was the gNMI SubscribeRequest operation, a dial-in mode. This has been, for a long time, one of the key benefits of gNMI versus NETCONF, which took more time to come up with its implementation definition. At the time of this writing, the adoption of streaming telemetry is much more mature in gNMI, which is widely adopted by most vendors. NETCONF dial-in and dial-out are still in their early stages.

Dial-out telemetry has some incipient implementations—for instance, by Cisco and Juniper, using gRPC as the transport protocol. In parallel, the IETF is working on standards grouped as YANG Push, UDP and HTTPS/TCP transport options, and support for JSON and CBOR encoding.

NOTE
Concise Binary Object Representation or CBOR standardized in RFC 8949, is a binary data format based on JSON that supports schema definition directly with the YANG language.

Model-driven telemetry has prioritized TCP as a transport protocol over UDP (used by SNMP) for providing more reliable data transfer, and with nonrepudiation. However, UDP support provides some benefits in highly intense event traffic, such as in sampling mode when TCP benefits are not mandatory.

As we will explore in Chapter 14 when discussing the role of model-driven telemetry in a network automation strategy, telemetry is usually combined with message brokers to distribute data from the collectors, adding some features such as data schema validation, versioning, or routing.

Which solution will prevail is hard to predict. gNMI is well established and supported by a lot of vendors. On the other side, YANG Push supports more use cases (e.g., more encoding options) and comes from a standardization body (IETF) that is important for some network industry actors (e.g., service providers).

Now that we’ve explained the API interfaces available to manage network devices and controllers, you must understand how to automate them via these APIs. We’ll now take a look at using Python and Go to automate these interfaces, and also SSH.

Using Network APIs for Automation
As we’ve stated, the tools for exploring and learning to use an API differ from the tools used to consume an API within a programmatic solution. Thus far, we’ve looked at cURL for exploring HTTP-based APIs, an interactive NETCONF over SSH session for exploring the use of NETCONF, and gNMIc for exploring gNMI. In this part of the chapter, we’ll look at how to use Python and Go to automate network devices, using some popular libraries:

Python Requests
An intuitive and popular HTTP library for Python. This is the library we will use for automating devices and controllers with both RESTful HTTP-based APIs and non-RESTful HTTP-based APIs.

Go net/http
A built-in package to serve as an HTTP client or server. Introduced in Chapter 7, it’s similar to Python Requests, and we will use it to demonstrate how to interact with RESTCONF.

Python ncclient
This is a NETCONF client for Python, so we will use it for automating devices using NETCONF.

Go OpenConfig gNMIc
We used this gNMI client as a CLI in the previous section. Here, we will use the package directly in Go applications to interact via the gNMI interface.

Python Netmiko
This is a network-first SSH client for Python. This is the library we will use for automating devices via native SSH for devices without programmatic APIs.

WARNING
Even though we cover multiple APIs in this chapter, it is meant to be read from start to finish and not as API documentation for any given API. All the scripts created in this chapter have loose error handling because we are prioritizing simplicity to show the ideas instead of creating production-ready code. Be aware that the examples shown depend on the library’s version. The syntax and signature can change from one version to another.

Let’s start by looking at the Requests library and communicating with HTTP-based APIs.

The Python Requests Library
You’ve seen how to make HTTP-based API calls from the command line with cURL, or maybe you used the Postman GUI. These are great mechanisms for learning how to use a given API—but realistically, to write a script or a program that helps automate network devices, you need to be able to make API calls from within a script or program. In this section, we introduce the Python Requests library, which simplifies working with web-based APIs.

To enable an easy mapping between the previous section and this one, we reuse the same examples from “Exploring HTTP-based APIs with cURL”, the Cisco Meraki API, and the Arista eAPI. This section is meant to be read from start to finish, as the core focus is getting started with using the Requests library.

TIP
To install Requests, you can use pip3 within your virtual environment. Remember that you can review basic Python concepts in Chapter 6. Here’s the installation command:

$ pip3 install requests
$ pip3 list | grep requests
requests                          2.26.0
Automating the Meraki API with Requests
Let’s dive in and take a look at our first example using Requests. We’re going to create a complete Python script to retrieve the first network, from the first organization, from a Cisco Meraki account. We’ve already executed this same GET request with cURL in Examples 10-2 and 10-3. Remember to check the Meraki API documentation if needed.

Now, in Example 10-22, we focus on the first part of a script to retrieve the organizations available in Cisco Meraki.

Example 10-22. Using Requests to get Meraki organizations
#!/usr/bin/env python3

# The Python Requests library is used to issue and work HTTP-based systems.
import requests

# This executes if our script is being run directly.
if __name__ == "__main__":
    # This token is taken from Cisco Developer Hub to experiment with the API
    my_token = "6bec40cf957de430a6f1f2baa056b99a4fac9ea0"

    # This statement creates a Python dictionary for the HTTP request
    # headers that are going to use in the API calls.  The first two
    # headers we are setting are Content-Type and Accept.
    # The last one uses a custom Meraki header for authentication
    headers = {
        "Accept": "application/json",
        "Content-Type": "application/json",
        "X-Cisco-Meraki-API-Key": my_token,
    }

    # The URL is saved as a variable called base_url to modularize our
    # code and simplify the next statement.
    base_url = "https://api.meraki.com/api/v1"

    # In the Requests library, there is a function per HTTP verb, and in
    # this example we are issuing a GET request, so we are therefore
    # using the get function. We pass two objects into the get function.
    # The first  object passed in must be the URL, and the others should be
    # keyword arguments (key=value pairs). Then, we pass the proper headers.
    response = requests.get(f"{base_url}/organizations", headers=headers)
At the end of this Python script, the response variable contains the response from the Meraki API. You could reproduce the same steps in an interactive Python session, as noted in Chapter 6. Actually, it’s a great idea to try this to better understand the content of the response variable. So, let’s run the previous script in the interactive interpreter, using the -i flag. It is a great way to test and troubleshoot:

ch10-apis/python_requests$ python3 -i get_networks.py
# The interactive execution leave us at the end of our script, after the
# response = requests.get(f"{base_url}/organizations", headers=headers)
>>> response
<Response [200]>
The response variable contains a 200 HTTP response, so you should assume everything worked as expected. Indeed, one important piece of data is missing: where is the content of the response, containing the organizations?

Being in the Python interactive interpreter session, you can inspect the response with dir(), displaying all attributes and methods of a given object:

>>> dir(response)
[ # output omitted for brevity
 'apparent_encoding', 'close', 'connection', 'content',
 'cookies', 'elapsed', 'encoding', 'headers', 'history',
 'json', 'links', 'next', 'ok', 'raise_for_status', 'raw', 'reason',
 'request', 'status_code', 'text', 'url']
From all the attributes and methods available in the response, we will review two: status_code and json.

The status_code attribute gives us access to the HTTP response code as an integer:

>>> print(response.status_code)
200
The json() method returns the response content as a dict, decoding the JSON contained in the text attribute (which stores the actual response):

>>> response.json()
[
  {
    'id': '573083052582915028',
    'name': 'Next Meraki Org',
    'url': 'https://n18.meraki.com/o/PoiDucs/manage/organization/overview',
    'api': {
      'enabled': True
    },
    'licensing': {
      'model': 'co-term'
    },
    'cloud': {
      'region': {
        'name': 'North America'
      }
    }
  },
  # other entries omitted for brevity
]
>>>
You can start using this output by saving it as a variable and extracting the value of the name key from the first organization:

>>> organizations = response.json()
>>> type(organizations)
<type 'list'>
>>> organizations[0]["name"]
'Next Meraki Org'
In our script, we also store the response data in an organizations variable, and you can resume the script that will get the networks from a given organization, as shown in Example 10-23.

Example 10-23. Using Requests to get Meraki networks
    # Continues from the previous code snippet

    # We pick the id from the first organization to later gather the related networks
    first_organization_id = organizations[0]["id"] 1

    # Similar get request, composing the URL with the organization id
    response = requests.get(
        f"{base_url}/organizations/{first_organization_id}/networks", headers=headers
    )

    networks = response.json()                     2
1
We pick the identifier from the first available organization, accessing the 0 index in the list of organizations.

2
If we explore the content of networks, we will see another dictionary—this time containing the network data from the Meraki API.

Let’s continue to build on this; you are now going to create a new network (within an organization) by using Requests, as in Example 10-4. As you may guess, you need to only extend the previous script, adding the same content used in the cURL example: update the URL, update the HTTP request type, and send data in the body of the request, as we do in Example 10-24.

Example 10-24. Using Requests to create a Meraki network
    # Continues from the previous code snippet
    # first_organization_id comes from the previous script

    # json library is used to encode a dictionary object as a JSON string
    import json

    # Payload contains the data necessary to define the expected
    # data to create a new network in the API
    payload = {
        "name": "my brand new automated network",
        "productTypes": ["switch"],
    }

    # Using the post method instead of get to create an object
    response = requests.post( 1
        f"{base_url}/organizations/{first_organization_id}/networks",
        headers=headers,
        data=json.dumps(payload)
    )

    print(response.json())
1
Pay attention to the HTTP verb being used. This particular request uses the post() function as a resource is being created. To update a resource, you’d use the patch() function, and to replace a resource, you’d use put(). You’ll see these functions in more examples later in the chapter.

Let’s focus on how to send data in the body of the HTTP request. This is where you need to differentiate between a Python dictionary and a JSON string. While you work with dictionaries in Python to construct the required body, this is sent over the wire as a JSON string. To convert the dictionary to a well-formed JSON string, you use the dumps() function from the json module. This function takes a dictionary and converts it to a JSON string. You finally take the string object and pass it over the wire by assigning it to the data key being passed to the post() function.

TIP
We typically know which payload data is required to create/update an API resource by checking the API documentation. Some of the attributes will use some defaults, but for the rest, we need to provide some data. In the example, name looks like an obvious mandatory piece of data. But productTypes could be missed, assuming that it could be a default product type. In these cases, the API response should provide a useful error message to guide you, pointing out what is missing in the request payload.

Finally, when running the script, you get a response (from the POST). It provides a new network with the name and the productTypes defined before, plus the rest of the parameters automatically assigned in creation:

ch10-apis/python_requests$ python3 create_network.py
{
  'id': 'N_573083052583238204',
  'organizationId': '573083052582915028',
  'productTypes': ['switch'],
  'url': 'https://n18.meraki.com/my-brand-new-aut/n/vAQKbcs/manage/usage/list',
  'name': 'my brand new automated network',
  'timeZone': 'America/Los_Angeles',
  'enrollmentString': None,
  'tags': [],
  'notes': None,
  'isBoundToConfigTemplate': False
}
You have learned about using the Python Requests library for the Cisco Meraki API. Next, we continue exploring Requests but for the Arista eAPI (as in the cURL examples).

Consuming eAPI in a Python script
We’re now going to look at Arista’s eAPI. As you learned in “Understanding non-RESTful HTTP-based APIs”, as you go through the next examples with eAPI, keep the following points in mind:

eAPI is a non-RESTful HTTP-based API. In other words, it’s an HTTP-based API that doesn’t follow all the principles of REST. An HTTP POST is used no matter which operation is being performed—​even if show commands are used, a POST is still used. Specifically, it uses the JSON-RPC protocol to communicate between you (the client) and the switch (the server).

Remember POST requests require data to be sent in the data payload of the request. This is where solid API tools and documentation come into play.

The URL format for eAPI API calls is always http(s)://<ip-address-eos>/command-api.

TIP
Arista switches have a built-in tool called the Command Explorer that you could leverage to learn the required structure of the payload object. The API documentation provides details about this tool.

We start with a basic Python script in Example 10-25. This code uses the Requests library to communicate to the Arista API (eAPI), executing the show vlan brief CLI command. This command should return the VLAN information with the device. At the end of the script, the code will output the response and the HTTP status code.

Example 10-25. Using Python Requests with Arista eAPI
import json
import sys
import requests
from requests.auth import HTTPBasicAuth

if __name__ == "__main__":

    # requests class to create the basic authentication header
    auth = HTTPBasicAuth('ntc', 'ntc123') 1

    url = 'http://eos-spine1/command-api' 2

    # payload expected by the API
    payload = {
        "jsonrpc": "2.0",
        "method": "runCmds",
        "params": {
            "format": "json",
            "timestamps": False,
            "cmds": [
                "show vlan brief"
            ],
            "version": 1
        },
        "id": "EapiExplorer-1"
    }

    # Even though we retrieve data with the "show vlan brief",
    # the API uses the POST method
    response = requests.post(url, data=json.dumps(payload), auth=auth) 3

    # Helper output onscreen to show the status code, and the response
    print(f'STATUS CODE: {response.status_code}')
    print(f'RESPONSE: {json.dumps(response.json(), indent=4)}')
1
HTTPBasicAuth class, from the Requests library, to create the basic authentication format (using a Base64 encoding).

2
Using an http endpoint is not recommended for production. If we use a self-signed certificate or unverified HTTPS connection (adding verify=False to the requests method), we will receive a warning, which can be disabled with requests.packages.urllib3.disable_warnings().

3
Even though we are retrieving data, the HTTP method used is POST. So, it requires a payload, which defines several parameters required by the API.

When you execute the script, it should give a similar output:

ch10-apis/python_requests$ python3 eapi-requests.py
STATUS CODE: 200
RESPONSE:
{               1
    "jsonrpc": "2.0",
    "id": "EapiExplorer-1",
    "result": [ 2
        {
            "sourceDetail": "",
            "vlans": {
                "1": {
                    "status": "active",
                    "name": "default",
                    "interfaces": {
                        "Ethernet1": {
                            "privatePromoted": false
                        },
                        # Omitted other interfaces for brevity
                    },
                    "dynamic": false
                },
                "20": {
                    "status": "active",
                    "name": "VLAN0020",
                    "interfaces": {},
                    "dynamic": false
                },
                "30": {
                    "status": "active",
                    "name": "VLAN0030",
                    "interfaces": {},
                    "dynamic": false
                }
            }
        }
    ]
}
1
The response is a nested JSON object, as expected.

2
The output of the commands (only one in this example) is a list of dictionaries. There would be a list element for each command executed.

It’s noticeable that, even though you send a CLI command (show vlan brief), the format of the response content is JSON. This makes it much easier to interact with the API programmatically. However, if you are interested in the traditional text CLI output, you can specify it with "format": "text" in the payload, and the response would contain an output key with a string:

RESPONSE: {
    "jsonrpc": "2.0",
    "id": "EapiExplorer-1",
    "result": [
        {
            "output": "VLAN  Name                             Status    Ports\n
            ---- -------------------------------- --------- -------------------"
            # output omitted for brevity
        }
    ]
}
Next, we go a bit further with a more elaborate example to show the potential of automating the network via API-based scripts.

Using eAPI to autoconfigure interface descriptions based on LLDP data
Let’s continue to use eAPI to build something a little more useful. How about a Python script that autoconfigures interface descriptions for Ethernet interfaces based on LLDP neighbors for two Arista spine switches?

To do this, you should modularize the script to support multiple devices as well as have a simple way to send multiple API calls without requiring multiple payload objects in the script. Our goal is to autoconfigure interface descriptions such that they will look like the following (this is an example, and you will see the actual LLDP data and the final description later):

interface Ethernet2
  description Connects to interface Ethernet2 on neighbor eos-leaf1.ntc.com
  no switchport
!
interface Ethernet3
  description Connects to interface Ethernet2 on neighbor eos-leaf2.ntc.com
  no switchport
!
To easily digest the code, we split it into three parts so you can progressively understand the complete example.

First, in Example 10-26, you create the issue_request() function, which takes two arguments: the target device and the commands. The data is the only data in the requests’ operation. So, with this helper function, you could later pass different target devices and different commands to obtain the response (already converted from JSON to a Python object). This is a good example of the Don’t Repeat Yourself (DRY) software development principle.

Example 10-26. Wrapping requests in a helper function
import json
import sys
import requests
from requests.auth import HTTPBasicAuth

# Helper method to issue "commands" to a "device", and return the result
def issue_request(device, commands): 1
    """Make API request to EOS device returning JSON response."""
    payload = {
        "jsonrpc": "2.0",
        "method": "runCmds",
        "params": {
            "format": "json",
            "timestamps": False,
            "cmds": commands,
            "version": 1
        },
        "id": "EapiExplorer-1"
    }

    response = requests.post(
      'http://{}/command-api'.format(device),
      data=json.dumps(payload),
      auth=HTTPBasicAuth('ntc', 'ntc123')
    )

    return response.json()

# continues in the next example
1
The code of this function is exactly the same as Example 10-25, but has been modularized for reusability.

Next, in Example 10-27, we leverage the issue_request() function to get the specific information we want from the API response (in this case, lldpNeighbors). This implies knowledge of the data structure from the response, which you can get by experience or from the documentation.

Example 10-27. Extracting LLDP neighbors from the response
# continues from the previous example

def get_lldp_neighbors(device):
    """Get list of neighbors

    Sample response for a single neighbor:
        {
          "ttl": 120,
          "neighborDevice": "eos-spine2.ntc.com",
          "neighborPort": "Ethernet2",
          "port": "Ethernet2"
        }
    """
    # Define the target methods
    commands = ['show lldp neighbors']
    response = issue_request(device, commands)
    # Extract the neighbors' data from the result of the first and only command
    # and return it as a list of dictionaries
    return response['result'][0]['lldpNeighbors']

# continues in the next example
TIP
Creating readable code takes practice, but in Example 10-27 you can observe two useful approaches:

Using self-descriptive naming: get_lldp_neighbors clearly defines its intent

Leveraging function docstrings to explain the function’s purpose—in this case, the format of the response

Finally, in Example 10-28, we add another helper function, configure_interfaces(), and the main function to run the script. The configure_interfaces() function does exactly what it describes: takes the list of neighbors, and with configuration commands, updates the description of the interfaces. In the main function, you define all the target devices to iterate on and perform two operations: get the LLDP information, and configure the interfaces description accordingly.

Example 10-28. Configuring the interfaces description with LLDP information
# continues from the previous example

def configure_interfaces(device, neighbors):
    """Configure interfaces in a single API call per device."""
    command_list = ['enable', 'configure']
    for neighbor in neighbors:
        local_interface = neighbor['port']
        if local_interface.startswith('Eth'):
            # Excluding Management as it has multiple neighbors
            description = (
              f"Connects to interface {neighbor['neighborPort']} on neighbor "
              f"{neighbor['neighborDevice']}"
            )
            description = 'description ' + description
            interface = f'interface {local_interface}'
            # Extending the list of commands, in the proper order
            command_list.extend([interface, description])
    # Retrieve the output from the commands created from the neighbors
    response = issue_request(device, command_list)

if __name__ == "__main__":
    # device names are FQDNs
    devices = ['eos-spine1', 'eos-spine2']
    for device in devices:
        neighbors = get_lldp_neighbors(device)
        configure_interfaces(device, neighbors)
        print('Auto-configured Interfaces for {}'.format(device))
NOTE
Going through this example, you may be wondering if you could have organized the code differently. Maybe you could combine get_lldp_neighbors with the configure_interfaces() function, to get a bigger one. Or you could call issue_request out of the other functions, in the main code. The point here is that you have a myriad of options to create valid code. Choose one, experiment with it, and look for better patterns toward reusability and readability, while keeping it simple.

Let’s run the script that will update the interfaces’ descriptions according to the LLDP neighbor:

ch10-apis/python_requests$ python3 eapi-autoconfigure-lldp.py
Auto-configured Interfaces for eos-spine1
Auto-configured Interfaces for eos-spine2
Using the Requests Python library is an easy way to interact with APIs in your Python applications. However, to make it even simpler, some APIs provide their own libraries, the SDKs. We will give a quick glance at SDKs next.

Using API SDKs
An API SDK is a software package that abstracts access to an API by using functions, methods, and/or classes. It allows faster development because it comes with all the common conventions implemented. Therefore, you don’t need to reinvent the wheel every time, reducing development time. The SDK makes the code simpler and more readable. On the other hand, it could introduce some constraints due to non-implemented features available in the API or introduce library dependencies, increasing the footprint of your application.

Most API platforms offer SDKs in the most popular languages in their user community. Both APIs explored in the previous section offer Python SDKs:

Cisco Meraki

Arista eAPI

It’s not the purpose of this book to document using any specific SDK, but showing how an SDK is used enables you to see what one looks like. For detailed information about an SDK, check its reference docs page.

Exploring the Meraki API SDK
Using the Meraki API SDK instead of the Requests library, you will get the same output as in Example 10-22 without having to know about some API conventions. For instance, you don’t need to know about the custom authentication header the API expects. This is also useful for maintainability because if this header key changes, you don’t need to update your code.

NOTE
Use pip3 to install the Meraki API SDK:

$ pip3 install meraki
$ pip3 list | grep meraki
meraki                            1.25.0
After importing the library, you have to instantiate the class meraki.DashboardAPI, which contains all the methods to interact with the API. This initialization requires only the API key used before, but not the URL or the authentication header key, as these are implicitly defined by the library:

>>> import meraki
>>>
>>> meraki_client = meraki.DashboardAPI(
...   api_key="6bec40cf957de430a6f1f2baa056b99a4fac9ea0")
2022-10-01 16:18:28 meraki: INFO > Meraki dashboard API session initialized with ...
# output omitted for brevity
>>>
Then you can retrieve the organization, as you did in Example 10-22. Instead of crafting the HTTP request, you use the getOrganizations() method in meraki_​cli⁠ent.organizations:

>>> my_orgs = meraki_client.organizations.getOrganizations()
2022-10-01 16:18:49 meraki: INFO > GET https://api.meraki.com/api/v1/organizations
2022-10-01 16:18:50 meraki: INFO > GET https://n392.meraki.com/api/v1/organizations
2022-10-01 16:18:50 meraki: INFO > organizations, getOrganizations - 200 OK
>>>
>>> my_orgs[0]
{
  'id': '573083052582915028',
  'name': 'Next Meraki Org',
  'url': 'https://n18.meraki.com/o/PoiDucs/manage/organization/overview',
  'api': {'enabled': True},
  'licensing': {'model': 'co-term'},
  'cloud': {'region': {'name': 'North America'}}
}
>>>
By now, you’re likely feeling comfortable interacting with HTTP APIs with Python. You’ve interacted with a native RESTful HTTP API, Cisco Meraki, and a non-RESTful one, Arista eAPI. As a reminder, every request to eAPI is an HTTP POST, and the URL is the same for every request, whereas a truly RESTful API using HTTP as its transport has a different URL based on the resource in question (e.g., organization, network, or routes in Cisco Meraki API).

Next, we will use the other programming language covered in the book, Go, to explore a specific HTTP API, the RESTCONF interface.

The Go net/http Package
In the previous section, you learned how to use the Python Requests library to interact with HTTP APIs. Here, we will use the Go net/http package to interact with HTTP APIs (focusing specifically on the RESTCONF interface—in this case, exposed by a Cisco IOS XE). Chapter 7 introduced the net/http package, and here we’ll show a few examples of interacting with the RESTCONF interface.

Using net/http with RESTCONF
First, let’s start with the same operation as in “Exploring RESTCONF in Cisco IOS XE” to obtain the full native configuration. In Example 10-29, the main differences from the net/http example in Chapter 7 are the use of the basic http.NewRequest to create a GET request and the Authorization header.

Example 10-29. Using Go net/http to retrieve the configuration
package main

import (
    "crypto/tls"
    "encoding/base64"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
)

// helper method to construct the expected format of
// the basic authentication string, encoded in base64
func basicAuth(username, password string) string {
    auth := username + ":" + password
    return base64.StdEncoding.EncodeToString([]byte(auth))
}

// helper method to implement the error-checking pattern
func checkError(err error) {
    if err != nil {
        log.Fatal(err)
    }
}

func main() {
    transCfg := &http.Transport{
        // ignore expired SSL certificates
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }
    // create a new HTTP client, with the previously defined transport config
    client := &http.Client{Transport: transCfg}

    // create a new HTTP request, with the method, url, and headers
    request, err := http.NewRequest("GET",
        "https://csr1/restconf/data/Cisco-IOS-XE-native:native", nil)
    checkError(err)
    request.Header.Set("Accept", "application/yang-data+json")
    request.Header.Add("Authorization", "Basic "+basicAuth("ntc", "ntc123"))

    // perform the HTTP request, defined before, and store it in result
    result, err := client.Do(request)
    checkError(err)
    // read the body content from the response
    body, err := ioutil.ReadAll(result.Body)
    checkError(err)
    result.Body.Close()
    fmt.Printf("%s", body)
}
Running the code (go run get_config.go) will return the very same output obtained in “Exploring RESTCONF in Cisco IOS XE”. However, you could narrow the output by extending the path in the NewRequest with a specific data model and including filtering.

For instance, you can retrieve a specific interface configuration by using URI-encoded path expressions. In get_config_interface_g1.go, we use the URI-encoded path expression interface=GigabitEthernet1 to retrieve the configuration of the GigabitEthernet1 interface: restconf/data/ietf-interfaces:interfaces/interface=GigabitEthernet1. This matches the identifier field (name) of the interface in the YANG model and gets you only the configuration of the specified interface:

ch10-apis/go_http$ go run get_config_interface_g1.go
{
  "ietf-interfaces:interface": {
    "name": "GigabitEthernet1",
    "description": "MANAGEMENT_DO_NOT_CHANGE",
    "type": "iana-if-type:ethernetCsmacd",
    "enabled": true,
    "ietf-ip:ipv4": {
      "address": [
        {
          "ip": "10.0.0.15",
          "netmask": "255.255.255.0"
        }
      ]
    },
    "ietf-ip:ipv6": {
    }
  }
}
Updating configuration via RESTCONF with net/http
You can also use net/http to run HTTP requests to change the configuration state. In “Exploring RESTCONF in Cisco IOS XE”, we used a PATCH request to update the existing OSPF configuration, and we mentioned that with a PUT request, you could apply the declarative configuration approach that defines the final state without considering the current one.

In Example 10-30, you can see the most relevant code from update_ospf_config.go that replaces the OSPF configuration with a new one. The HTTP request takes a JSON payload with the expected data structure and uses it in the PUT request.

Example 10-30. Using Go net/http to update the OSPF config
// Omitted code

func main() {
  // Omitted code

  // JSON payload for HTTP request
  var jsonStr = []byte(`{
  "router-ospf": {
    "ospf": {
      "process-id": [
        {
          "id": 10,
          "network": [
            {
              "ip": "203.0.113.0",
              "wildcard": "0.0.0.7",
              "area": 0
            },
            {
              "ip": "203.0.113.64",
              "wildcard": "0.0.0.7",
              "area": 0
            }
          ],
          "router-id": "203.0.113.1"
        }
      ]
    }
  }
}'`)

  // create a new HTTP request, with PUT method, new URL,
  // the payload, and new headers
  request, err := http.NewRequest("PUT",
    "https://csr1/restconf/data/Cisco-IOS-XE-native:native/router/
    Cisco-IOS-XE-ospf:router-ospf", bytes.NewBuffer(jsonStr))
  checkError(err)
  request.Header.Set("Accept", "application/yang-data+json")
  request.Header.Set("Authorization","Basic " + basicAuth("ntc","ntc123"))
  request.Header.Set("Content-Type", "application/yang-data+json")

  // Omitted code
}
NOTE
In a real example, you would get the payload from an external file or a database, but we used a JSON object directly to simplify the script.

After these changes, you run the script, and voilà! The new OSPF configuration is in place. You can check it by repeating the GET request to the OSPF configuration node in the router.

With RESTful APIs that offer this type of power and control, you need to ensure you have a good process for making changes. As you can see, if you are trying to make only a small change or addition and happen to send a PUT request, catastrophic consequences can result. From an overall adoption perspective of this particular API, you may want to start using PATCH requests and gradually migrate to the point where you can indeed use PUT to declaratively manage specific sections of the configuration.

WARNING
Example 10-30, using the PUT declaratively, configures everything under the ospf key. But if you targeted the router key, this key change would technically eliminate all other routing protocol configurations too. That’s not a good thing. However, we chose to show the power and potential danger if it’s not used and understood properly.

Now that we’ve shown how to start automating devices that have HTTP-based APIs by using Python and Go, let’s look at the same approach by using the Python ncclient for automating devices with the NETCONF interface.

The Python ncclient Library
The Python ncclient is a popular NETCONF client for Python. It is client software that is built to communicate programmatically with NETCONF servers. Remember, in our case, a NETCONF server is going to be a network device. We’ll walk through a vMX Juniper example, but the same approach will work for other platforms.

NOTE
To install ncclient, you can use pip3:

$ pip3 install ncclient
$ pip3 list | grep ncclient
ncclient                          0.6.13
Once ncclient is installed, you can start to issue NETCONF API calls to network devices. Let’s walk through this while in the Python interactive interpreter.

NOTE
ncclient is not the only Python library offering NETCONF capabilities. Vendor-specific ones exist, such as PyEZ for Junos, as well as other generalist ones such as scrapli-netconf and netconf-client.

When you enter the Python interpreter, your first step is to import the manager module within the ncclient Python package:

>>> from ncclient import manager
The basic function we are going to use within the manager module is responsible for establishing a persistent connection to the device. Keep in mind that since NETCONF runs over SSH, this connection is stateful and persistent (as compared to RESTful APIs being stateless). This function is called connect() and accepts several parameters such as hostname/IP address, port number, and credentials. You’ll see in the following example that other parameters will remain unchanged that map back to the underlying SSH configuration and properties:

>>> device = manager.connect(
...   host='vmx1', port=830, username='ntc',
...   password='ntc123', hostkey_verify=False,
... )
As soon as the connect() function is called, a NETCONF session is established to the network device, and an object is returned and saved as device, an instance of a ncclient Manager object. This object exposes the methods to interact with NETCONF RPC operations.

Understanding the Manager object
First, let’s look at the device capabilities, as in Example 10-8, and get exactly the same output, but in a different format (a Python list):

>>> list(device.client_capabilities)
[
  'urn:ietf:params:netconf:base:1.0',
  'urn:ietf:params:netconf:base:1.1',
  'urn:ietf:params:netconf:capability:writable-running:1.0',
  'urn:ietf:params:netconf:capability:candidate:1.0',
  'urn:ietf:params:netconf:capability:confirmed-commit:1.0',
  # output omitted for brevity
]
Retrieving Juniper vMX device configurations with ncclient
In Example 10-9, when exploring the use of NETCONF over an SSH session, you used the complete <get> XML RPC operation to query the vMX router for the fxp0 interface configuration. Now, you reuse the subtree filtering, defined as a string:

>>> get_filter = """
...   <configuration xmlns="urn:ietf:params:xml:ns:netconf:base:1.0">
...     <interfaces>
...       <interface>
...         <name>fxp0</name>
...       </interface>
...     </interfaces>
...   </configuration>
... """
NOTE
Remember that triple quotes in Python denote a multiline comment and can be used to create a multiline string that can be used as a value of a variable.

Once the filter is defined, you pass that as a parameter to the get() method, specifying the subfilter subtree type. Without the filter, you would get the full configuration:

>>> nc_get_reply = device.get(('subtree', get_filter))
At this point, after the NETCONF GET request to the device, the result is stored in nc_get_reply and matches the result in Example 10-10.

NOTE
The ncclient supports XPath filters. But, this support depends on the NETCONF capabilities of the platform. All filters used with Junos vMX examples are subtree filters because this platform doesn’t support XPath for get operations. You can see how to use XPath filters in “Using ncclient with Cisco IOS XE”.

While our examples use XML strings as the filters, it is also possible to use native XML objects (etree objects). We are using string objects because they are much more human-readable and easier to use when getting started. You may want to use native etree objects if you need to dynamically build a filter object.

The data attribute in nc_get_reply contains the native XML object from the lxml Python library (lxml is covered in more detail in Chapter 8). For instance, the etree.tostring() method converts a native XML object to a string:

>>> from lxml import etree
>>>
>>> as_string = etree.tostring(nc_get_reply.data)
>>> print(as_string)
b'<rpc-reply message-id="urn:uuid:e2c1daa0-8556-4e6b-84dc-e72e90809f73"><data>
<configuration commit-seconds="1653021086" commit-localtime="2022-05-20 04:31:26
UTC" commit-user="ntc"><interfaces><interface><name>fxp0</name><unit><name>0
</name><description>MANAGEMENT_INTERFACE__DO_NOT_CHANGE</description><family>
<inet><address><name>10.0.0.15/24</name></address></inet></family></unit>
</interface></interfaces></configuration><database-status-information>\n
</database-status-information></data></rpc-reply>'
In the output of nc_get_reply, you can observe the IP address and mask for the fxp0 interface (10.0.0.15/24). If this is the information you are interested in, you should parse the XML tags that contain this information: <address> and <name>.

To find a specific tag in an lxml.etree object, you can use the find() method. It provides a simple way to search a full XML object for a given XML tag when using the expression denoted by .//. Since you want to extract the <address> object and its children, you could try the following example:

>>> address = nc_get_reply.data.find('.//address')
Unfortunately, it doesn’t work. What’s wrong here? The statement tries to extract the XML element with the <address> tag. However, when XML namespaces are used, the tag name must include the namespace concatenated—in other words, {namespace}tag. In this case, {http://yang.juniper.net/junos/conf/interfaces}address. Alternatively, if an XML namespace alias is defined, it can be used as alias : tag.

NOTE
Our example has multiple namespaces. You can gradually print one child object at a time to see which namespace is used. In the example, the default namespace is urn:ietf:params:xml:ns:netconf:base:1.0, but when you print a single object, you see only one. The next namespace in the hierarchy is overriding the default namespace for all children of the <configuration> element.

To properly extract the IP address and mask, we’ll follow three steps. First, we extract the <address> object that contains an inner <name> XML tag, then we extract the content of the <name> tag with the prefix and mask combined, and finally, we save it as a string by using text:

>>> address = nc_get_reply.data.find(
>>>   './/{http://yang.juniper.net/junos/conf/interfaces}address'
>>> )
>>> ip_address = address.find(
>>>   './/{http://yang.juniper.net/junos/conf/interfaces}name' 1
>>> )
>>> ip_address.text
'10.0.0.15/24'
1
We can’t look for <name> directly because we would hit the interface name, which is first in the hierarchy.

You may be thinking, “Extracting values based on the namespaces is tedious.” You are absolutely right. Parsing XML is not easy, and it is even more difficult when namespaces are involved. But, in the end, once you know the namespace, you just need to concatenate two strings. Also, it’s possible to strip namespaces from an XML object before doing XML parsing, further simplifying the process.

By now, you should be getting the hang of issuing NETCONF <get> requests. Let’s look at one more, but this time, working with the SNMP configuration to use the findall() method.

Using findall() to retrieve multiple XML objects
On our Juniper vMX, we currently have two SNMP read-only community strings configured. For verification, this is the output after we issue the show snmp command while in configuration mode:

ntc@vmx1# show snmp
location EMEA-IE;
contact Scott_Grady;
community public123 {
    authorization read-only;
}
community supersecure {
    authorization read-write;
}
[edit]
TIP
Our goal is to extract the name of each community string and the authorization level for each. Juniper has functionality in its CLI such that you can see the expected XML response as well when you pipe the command to display xml (for instance, show snmp | display xml).

When you know the data returned from the NETCONF request, you can more easily write the associated Python code. Whenever you are issuing a NETCONF get request to a Junos device, <configuration> needs to be the outermost XML tag when you’re collecting configuration state information. Within that element, you can build the appropriate filter, which can be gleaned from the XML text found while on the CLI. The filter string to request SNMP configuration looks like this:

get_filter = """
... <configuration>
...   <snmp>
...   </snmp>
... </configuration>
... """
The next step, in Example 10-31, is to make the request, just as we’ve done already. After the request is made, we’ll verify the output, printing the XML string to the terminal.

Example 10-31. ncclient get SNMP configuration
>>> nc_get_reply = device.get(('subtree', get_filter))
>>> print(nc_get_reply)
<rpc-reply message-id="urn:uuid:c3170685-e275-4db1-855d-bb1a56404d55">
  <data>
    <configuration commit-seconds="1653021086"
    commit-localtime="2022-05-20 04:31:26 UTC" commit-user="ntc">
      <snmp>
        <location>EMEA-IE</location>
        <contact>Scott_Grady</contact>
        <community>
          <name>public123</name>
          <authorization>read-only</authorization>
        </community>
        # output omitted for brevity
  </data>
</rpc-reply>
As we stated, our goal is to parse the response, saving the community string and authorization type for each community. Rather than just print these to the terminal, let’s save them as a list of Python dictionaries. To do this, you’ll follow the same steps used earlier in Example 10-32.

Example 10-32. Using findall() to iterate over multiple XML objects
>>> snmp_list = []
>>> xmlns = "http://yang.juniper.net/junos/conf/snmp"
>>> communities = nc_get_reply.data.findall(f'.//{{{xmlns}}}community')    1
>>> for community in communities:                                          2
...     temp = {}
...     temp['name'] = community.find(f'.//{{{xmlns}}}name').text          3
...     temp['auth'] = community.find(f'.//{{{xmlns}}}authorization').text 3
...     snmp_list.append(temp)
...
>>>
>>> print(snmp_list)
[
  {'auth': 'read-only', 'name': 'public123'},
  {'auth': 'read-write', 'name': 'supersecure'}
]
1
Instead of using find(), youwe use the findall() method, which allows extracting multiple elements of the same type instead of only the first match. Using a string formation option, such as f-strings, simplifies xmlns management.

2
communities contains a list of objects, not only the first match.

3
Because we need to output the curly brackets in the output within the f-string, we need to escape them via doubling: {{{.

You’ve seen how to issue NETCONF requests to obtain configuration data, but now we are going to transition a bit and show how to make configuration changes via NETCONF by using the <edit-config> operation.

Making Junos vMX configuration changes with ncclient
To illustrate how to use ncclient to change the configuration via NETCONF, we will use the edit_config() method of the device object, which maps directly to the <edit-config> NETCONF operation. Following the previous example, now we’re going to configure a new SNMP community string.

The ncclient edit_config() method takes two mandatory arguments. The first one, called <target>, defines which configuration data store is going to get modified in the request. Valid data stores are running, startup, and candidate. The second parameter, called <config>, needs to be an XML string or object that defines the requested configuration changes.

Therefore, we need to construct the <config> that matches the expected format. The easiest way to get an understanding of the expected data structure is to use it from the get request, as shown in Example 10-31. Then, with this reference, we can construct a new XML configuration object. Let’s add a new community string called myNewCommunity that has read-only privileges:

>>> config_filter = """
... <config>
...   <configuration>
...     <snmp>
...       <community>
...         <name>myNewCommunity</name>
...         <authorization>read-only</authorization>
...       </community>
...     </snmp>
...   </configuration>
... </config>
... """
NOTE
You need to encapsulate this document in one final tag: <config>. This is often required when you’re using the <edit-config> operation, as covered in “Using NETCONF”.

To change the candidate data store with the new SNMP community, you call the edit_config() method:

>>> response = device.edit_config(target='candidate', config=config_filter)
>>> print(response)
<rpc-reply message-id="urn:uuid:56584f09-24d2-4e28-aa55-e583bf9d2ff2">
  <ok/>
</rpc-reply>
If you check the configuration via the CLI or make a get request, you’ll see a new community string on the device.

Performing NETCONF delete/replace operations with the ncclient
You’ve seen how to make a configuration change on the device by using the <edit-config> operation. The default edit operation is merge. However, more operations are available via the XML operation attribute.

After adding the SNMP community myNewCommunity, we are going to remove it by using the delete operation. In this case, we need to construct the new <config> with operation="delete" in the <community> tag. Notice that only the <name> tag is required to identify the object:

>>> config_filter = """
... <config>
...   <configuration>
...     <snmp>
...       <community operation="delete">
...         <name>myNewCommunity</name>
...       </community>
...     </snmp>
...   </configuration>
... </config>
... """
>>> response = device.edit_config(target='candidate', config=config_filter)
The merge and delete operations allow managing the configuration from its current state. But what if you want to implement a declarative approach, in which you don’t need to care about the current state, but only the final one?

There are a few ways to go about replacing a given hierarchy of XML configuration. One option is to make an API call to configure all desired SNMP community strings, and then retrieve all currently configured SNMP communities in another API call, loop over the response, and issue a delete operation per community for any not desired. While this approach is not terrible, NETCONF offers a better way—as you already know, of course.

With NETCONF, you can use the replace operation instead of merge or delete. This is the same as doing a PUT with RESTCONF. You define the replace operation (<snmp operation="replace">) and then define the desired SNMP community name and privilege.

When you start experimenting with the delete and replace operations, you need to be extremely careful. The merge operation is the default for a very good reason—​it enables you to only add or update a configuration. More importantly, if you put the operation="delete" line at the wrong place in the XML hierarchy, it could have a catastrophic effect on the network. For instance, placing the replace operation as an attribute in the <snmp> tag replaces the full configuration of SNMP. As we said with RESTCONF PUT, the same is true with NETCONF replace operations; be aware of their power.

WARNING
Always make sure to test in a lab or sandbox environment first. Do not test full replace operations in production!

Now, after exploring NETCONF automation with Juniper vMX, we want to give a glimpse into using it for Cisco IOS XE to show another implementation.

Using ncclient with Cisco IOS XE
To illustrate the usage of the XPath filter instead of the subfilter one, we will use Cisco IOS XE, which implements this capability. In this example, we want to extract the configuration for a specific interface (GigabitEthernet1).

>>> from ncclient import manager
>>> device = manager.connect(
>>>     host='csr1', port=830, username='ntc',
>>>     password='ntc123', hostkey_verify=False,
>>> )
>>> nc_get_reply = device.get(
...   ('xpath', '/interfaces/interface[name="GigabitEthernet1"]/config') 1
... )
1
The XPath filter injects [name="GigabitEthernet1"] in the path.

Now, as in Example 10-32, you use find() to traverse the XML structure to locate the data you are interested in. In this case, let’s extract the description of the interface GigabitEthernet1:

>>> description = nc_get_reply.data.find(
...   './/{http://openconfig.net/yang/interfaces}description')
>>> interface.text
'MANAGEMENT_DO_NOT_CHANGE'
As more vendors and operating systems support vendor-neutral data models, you’ll be able to issue the same exact API call against those devices, simplifying working with various device types. Until that point comes, you’ll need to understand the XML objects supported per NOS.

Understanding vendor-specific NETCONF operations
We have been focusing on the two most commonly used methods in ncclient—namely, <edit-config> and <get>. However, as you may have noticed earlier in this chapter, many more methods are available when working with NETCONF and ncclient. Here are a few of those:

commit()
Commits a candidate configuration to the active running configuration.

copy_config() and delete_config()
Creates or replaces, and deletes, respectively, an entire configuration data store with the contents of another complete configuration data store.

lock() and unlock()
For a production environment, you may want to lock the configuration data store before changes are made so that no other person or system can make changes during your NETCONF session. When it’s complete, you can unlock the configuration.

Everything we focused on with the Python-based ncclient was vendor neutral and would work across vendors, assuming you understand how to build the proper XML objects. However, you should understand that not every vendor implements NETCONF the same way even though it is an industry-standard protocol. For example, particular vendors have created their own platform-specific NETCONF operations or created their own methods for ncclient to simplify performing common operations. This is in contrast to NETCONF standard operations such as <edit-config>, <get>, <lock>, <unlock>, and <commit>.

To use these vendor-specific options when using ncclient, you need to specify the correct platform in the device_params parameter when instantiating a device object.

Every example we showed used device_params={} because we used industry-standard operations within ncclient. If you choose the vendor-specific methods and operations, you would set the device_params parameter to its required value. If you append device_params={"name": "junos"} in the connect() method, you will find out some custom Juniper RPC operations:

>>> device = manager.connect(
...  host='vmx1', port=830, username='ntc',password='ntc123',
...  hostkey_verify=False,device_params={"name": "junos"}
... )
>>> device._vendor_operations.keys()
dict_keys(['rpc', 'get_configuration', 'load_configuration', 'compare_configuration',
'command', 'reboot', 'halt', 'commit', 'rollback'])
NOTE
Juniper has developed custom methods within ncclient such as load_configuration(), get_configuration(), compare_configuration(), and command(), just to name a few. Several of Juniper’s methods are wrappers to simplify performing common tasks with standard NETCONF operations, and others use Juniper-specific NETCONF RPC operations.

For instance, using the commit() method, the candidate configuration can be moved to the running one, making the changes active:

>>> response = device.commit()
>>> print(response)
<rpc-reply message-id="urn:uuid:12a2bea3-7aa3-48dd-bfe8-569ccd84d61d">
  <ok/>
</rpc-reply>
After this complete overview of NETCONF (using the ncclient library), we will use the gNMIc Go package, the same we used in its CLI form before, to interact via the gNMI interface programmatically.

The OpenConfig gNMIc Go Package
In “Exploring gNMI with gNMIc”, you already used OpenConfig gNMIc to explore gNMI, but via CLI. Here, you will use the inner Go package to build Go scripts, emulating the same operations as before with an Arista EOS switch.

NOTE
In the Python world, the most popular library for gNMI management is pyGNMI, which offers a similar interaction with gNMI methods as ncclient offers for NETCONF. This is an example of performing a GetRequest to a specific path with the pyGNMI library:

from pygnmi.client import gNMIclient

with gNMIclient(
    target=('eos-spine1', '6030'), username='ntc',
    password='ntc123', insecure=True
  ) as gnmi_client:
    result = gnmi_client.get(path=[
      'openconfig-interfaces:interfaces', 'openconfig-acl:acl'])
As we did in “gNMI GetRequest”, we will start first with a GetRequest operation.

Using OpenConfig gNMIc to perform a gNMI Get operation
In this example, we create a basic programmatic script with gNMI. It contains the base program structure to interact with the gNMI interface that we will reuse in later examples. To extend your understanding of the gNMIc package, see the online documentation.

NOTE
To use the OpenConfig gNMIc package, you need to get and install the package first, or manage the Go module dependencies with go mod:

$ go get github.com/openconfig/gnmic/api
go: downloading github.com/openconfig/gnmic v0.27.1
go: added github.com/openconfig/gnmic v0.27.1
Example 10-33 defines the necessary imports. The most important ones are the openconfig/gnmic that offers access to the gNMI interface, and the protobuf that is used to manage the data serialization.

Example 10-33. gNMIc script, part 1
package main

import (
  "context"
  "fmt"
  "log"

  // gnmic package
  "github.com/openconfig/gnmic/api"

  // prototext marshals and unmarshals protocol buffer messages
  // as the textproto format, which offers generic support for text-based
  // request/response protocols
  "google.golang.org/protobuf/encoding/prototext"
)

// helper method to implement the error-checking pattern
func checkError(err error) {
  if err != nil {
    log.Fatal(err)
  }
}

// continues to part 2
Continuing the script, in Example 10-34, the five steps to retrieve the device configuration are defined.

Example 10-34. gNMIc script, part 2
// -> comes from previous

func main() {
  // create a gnmic target
  tg, err := api.NewTarget(
      api.Name("gnmi example"),
      api.Address("eos-spine1:6030"),
      api.Username("admin"),
      api.Password("admin"),
      api.Insecure(true),
  )
  checkError(err)

  // cancelable context releases the associated resources, as soon as the
  // operation is completed or canceled.
  ctx, cancel := context.WithCancel(context.Background())
  defer cancel()

  // create a new gNMI client within the base target, using the previous context
  err = tg.CreateGNMIClient(ctx)
  checkError(err)
  defer tg.Close()

  // create a GetRequest
  getReq, err := api.NewGetRequest(
      // Retrieve the full path, all the configuration
      api.Path("/"),
      // Define the expected payload encoding
      api.Encoding("json_ietf"))
  checkError(err)
  fmt.Println(prototext.Format(getReq))

  // send the created gNMI GetRequest to the gnmic target
  getResp, err := tg.Get(ctx, getReq)
  checkError(err)
  fmt.Println(prototext.Format(getResp))
}
NOTE
Notice the common Go pattern for error checking after each operation, wrapped in the checkError() function.

Then, it’s time to run the script and check its output:

ch10-apis/go_gnmic$ go run get_config.go
path: {}
encoding: JSON_IETF

notification:  {
  timestamp:  1664733418878706841
  update:  {
    path:  {}
    val:  {
      json_ietf_val:  "{\"openconfig-acl:acl\":{\"state\":
      {\"counter-capability\":\"AGGREGATE_ONLY\"}},
      \"arista-exp-eos:arista\":{\"eos\":
      {\"arista-exp-eos-igmpsnooping:bridging\":
      {\"igmpsnooping\":{\"config\":{}}},\"arista-exp-eos-mlag:mlag\":
      {\"config\":{\"dual-primary-action\":\"action-none\"
      "
      # omitted a VERY LONG output
    }
  }
}
In this initial example, we targeted the full data set with the / path. Now, as in Example 10-21, we narrow the scope to get the only interfaces’ data section—in particular, its configuration.

The only change from the previous script is changing the Path to filter the section /interfaces/interface/config, as defined in go_gnmic/get_interface_​con⁠fig.go:

  getReq, err := api.NewGetRequest(
    // Narrow the scope to the interfaces config
    api.Path("/interfaces/interface/config"),
    api.Encoding("json_ietf")
  )
If you run the script with the path change, you get a filtered output:

ch10-apis/go_gnmic$ go run get_interfaces_config.go
path:  # omitted path
encoding:  JSON_IETF

notification: {
  timestamp: 1664771809792522887
  update: {
    path: {
      # omitted path
    }
    val: {
      json_ietf_val: "{ 1
        \"openconfig-interfaces:description\":\"New Description\",
        \"openconfig-interfaces:mtu\":0,
        \"openconfig-interfaces:name\":\"Management0\",
        \"openconfig-interfaces:type\":\"iana-if-type:ethernetCsmacd\"}"
    }
  }
}
1
The interface’s configuration values, with the JSON_IETF encoding

The next step is to interact with the output of the Get operation, which we saved in the getResp variable.

Evolving from the previous script, in Example 10-35 (part of the example script get_interfaces_description.go), you narrow the path to the description leaf with /interfaces/interface/config/description. Then, you need only to extract its value and start using it.

Example 10-35. Using value from the gNMI response
  descriptionGnmiValue := getResp.GetNotification()[0].GetUpdate()[0].GetVal() 1
  myCurrentDescriptionValue, err := value.ToScalar(descriptionGnmiValue)       2
  checkError(err)

  myCurrentDescriptionStr := myCurrentDescriptionValue.(string)                3
  fmt.Println("This is my current description: " + myCurrentDescriptionStr)
1
Extracts json_ietf_val into the descriptionGnmiValue variable

2
Converts the variable value into a scalar value

3
Casts it into a string, and we are ready to use it

If you run the script, you can see how the interface description is retrieved and printed in the screen output:

ch10-apis/go_gnmic$ go run get_interfaces_description.go
# omitted previous output

This is my current description: New Description
With the current description, we will reuse part of the previous code to update the interface’s description with a new value.

Using OpenConfig gNMIc to perform a gNMI Set operation
As in “gNMI SetRequest” via CLI, in Example 10-36, you will update the interface description for the Management0 interface, but programmatically.

Example 10-36. gNMIc script to update configuration
  myNewDescription := myCurrentDescriptionStr + "_something_else" 1

  // create a gNMI SetRequest
  setReq, err := api.NewSetRequest(
      api.Update(                                                 2
          // Use XPath to target the description for Management0 interface
          api.Path("interfaces/interface[name=Management0]/config/description"),
          // Define the value to update
          api.Value(myNewDescription, "json_ietf")
      ),
  )
  checkError(err)

  // send the created gNMI SetRequest to the created target
  setResp, err := tg.Set(ctx, setReq)                             3
  checkError(err)
  fmt.Println(prototext.Format(setResp))
1
We use the previous description value, stored in myCurrentDescriptionStr, to create a new description.

2
The SetRequest type contains an Update object with the path and the value. gNMI supports two Set operation modes: update and replace (similar to a PUT versus PATCH in REST APIs).

3
We perform a Set operation in the gNMIc target, instead of a Get.

We run the script and receive confirmation that the update operation has been performed successfully:

ch10-apis/go_gnmic$ go run set_interfaces_description.go
# omitted previous output
response: {
  path: {
    # path omitted
  }
  op: UPDATE
}
timestamp: 1664863404149182568
If you retrieve the interface description configuration again, as before, you will see the new description updated.

Using OpenConfig gNMIc to subscribe to events
As we explained in “Understanding model-driven telemetry”, gNMI’s support for streaming telemetry via the subscribe operation was, and still is, one of the main drivers for its adoption. In Example 10-37, the goal is to subscribe to interface counters as you did in “gNMI Subscribe” via the CLI.

Example 10-37. gNMIc subscription to interface counters
  // create a gNMI subscribeRequest
  subReq, err := api.NewSubscribeRequest(
    api.Encoding("json_ietf"),
    // Select the stream mode, instead of poll
    api.SubscriptionListMode("stream"),
    // Define the subscription scope
    api.Subscription(                  1
      // Data container to retrieve
      api.Path("/interfaces/interface/state/counters"),
      // Define the subscription method. Others are on_change and target_defined
      api.SubscriptionMode("sample"),
      // For the sample mode, the sample interval is defined
      api.SampleInterval(10*time.Second),
    ))
  checkError(err)
  fmt.Println(prototext.Format(subReq))

  // start the subscription, identified as sub1, in a new goroutine
  go tg.Subscribe(ctx, subReq, "sub1") 2

  // start a goroutine that will stop the subscription after 30 seconds
  go func() {                          3
    select {
    case <-ctx.Done():
      return
    case <-time.After(30 * time.Second):
      // If the context is not stopped before, after 30 seconds
      // in this goroutine, it will stop the subscription defined
      // as sub1
      tg.StopSubscription("sub1")
    }
  }()

  // In the main process, starts creating two subscriptions, creating two channels
  // that are read in an infinite loop, until the subscriptions are closed
  // by the other goroutines
  subRspChan, subErrChan := tg.ReadSubscriptions()
  for {                                4
    select {
    case rsp := <-subRspChan:
      fmt.Println(prototext.Format(rsp.Response))
    case tgErr := <-subErrChan:
      log.Fatalf("subscription %q stopped: %v", tgErr.SubscriptionName, tgErr.Err)
    }
  }
1
SubscribeRequest is a sampling type, and it targets the state, not the config container.

2
The Subscribe() method is run in a different goroutine (notice the go before). This will keep it running while you progress in the script.

3
This goroutine is in charge to stop the subscription after 30 seconds. The identifier (sub1) connects the dots.

4
The main process will keep reading the subscription responses and errors (via two Go channels) until the subscription is stopped.

NOTE
Example 10-37 shows how three flows are executed and communicated to Go. The main one (initializing and printing the output received via channels) is the one managing the subscription, and the one to stop the subscription after a specified number of seconds.

When running the script, you will see the operation summary at the top, with the mode (SAMPLE), its interval (10 seconds), and the encoding. Remember that there is another mode type, on-change, which only sends data when there is a change in the data source. You can see the output here:

ch10-apis/go_gnmic$ go run subscribe_int_counters.go
subscribe:  {
  subscription:  {
    path:  {
      # path omitted
    }
    mode:  SAMPLE
    sample_interval:  10000000000
  }
  encoding:  JSON_IETF
}

update:  {
  timestamp:  1664944947620748134
  prefix:  {
    # prefix omitted
  }
  update:  {
    path:  {
      elem:  {
        name:  "out-octets"
      }
    }
    val:  {
      uint_val:  6243543
    }
  }
  update:  {
    path:  {
      elem:  {
        name:  "out-unicast-pkts"
      }
    }
    val:  {
      uint_val:  31339
    }
  }
  # output omitted for brevity
}
We’ve now provided an introduction to using the Python and Go libraries to communicate with modern programmatic network APIs. Now, we are going to shift gears and talk about using SSH in Python, as SSH is still the most widely deployed interface on network devices.

The Netmiko Python Library
Using CLI commands, SSH has been the de facto way network engineers and operators manage their infrastructure. Commands are passed over a persistent SSH connection to a network device, the device interprets them, and it responds with text that is viewable by a human on a terminal window. SSH does not use structured encoded data such as XML or JSON over the wire. While SSH is not a modern or a programmatic API, it is important to have an understanding of how to use Python to automate network operations with SSH for three reasons:

Not all devices support a programmatic API.

You may want to automate the turning on of the API.

Even if you’re automating a device with an API:

It’s good to have a backup plan.

Not all operations of a device may be supported with the API. This is not ideal, as it shows immaturity in the underlying API.

In this section, we show how to get started with a popular open source SSH client for Python called Netmiko. Netmiko’s purpose is to simplify SSH device management specifically for network devices.

NOTE
Do not underestimate the utility of the SSH Netmiko library, as it’s proven useful in providing a smooth transition from traditional network CLI management to network automation, and it’s heavily used and developed, as you can see in the Netmiko GitHub Contributors page.

We’re focused on Netmiko, as it provides a lower barrier to entry and already understands how to communicate with many network device types. Netmiko has varied support for dozens of device types, including those from Arista, Brocade, Cisco, Dell, HPE, Juniper, Palo Alto Networks, Linux, and many more; check the documentation for updated information. The great thing about Netmiko is that the overall usage is common across vendors. Only the commands used are specific to each platform.

NOTE
To install Netmiko, you can use pip3:

$ pip3 install netmiko
$ pip3 list | grep netmiko
netmiko                           3.4.0
The first thing you need to do is import the proper Netmiko device object. This object handles the SSH connection setup, teardown, and the sending of commands to the device. You used a similar approach with ncclient:

>>> from netmiko import ConnectHandler
You’re now ready to establish an SSH connection to the network device and create a Netmiko device object. The ConnectHandler object handles the SSH connection to the network device:

>>> device = ConnectHandler(
...   host='nxos-spine1',
...   username='admin',
...   password='admin',
...   device_type='cisco_nxos'
... )
At this point, there is an active SSH connection from Python using Netmiko with a Cisco NX-OS switch. Because each platform supports different commands and handles SSH differently, you must provide the device_type parameter when instantiating an instance of the ConnectHandler object.

Let’s check the available methods for our new device object called device by using the dir() function:

>>> dir(device)
[
  # methods removed for brevity
  'cleanup', 'clear_buffer', 'close_session_log', 'commit', 'config_mode',
  'conn_timeout', 'device_type', 'disable_paging', 'disconnect', 'enable',
  'encoding', 'establish_connection', 'exit_config_mode', 'exit_enable_mode',
  'select_delay_factor', 'send_command', 'send_command_expect',
  'send_command_timing', 'send_config_from_file', 'send_config_set',
]
As a network engineer, you should feel pretty comfortable with many of the attributes shown from the dir() function, as they are very network centric. We’ll walk through a few of them now.

Verifying the device prompt
Use the find_prompt() method to check the prompt string of the device:

>>> device.find_prompt()
'nxos-spine1#'
Entering configuration mode
Because Netmiko understands multiple vendors and what configuration mode means, it has a method to go into configuration mode that works across vendors; of course, the commands Netmiko uses under the covers may be different per OS:

>>> device.config_mode()
>>>
>>> device.find_prompt()
'nxos-spine1(config)#'
WARNING
Some NOSs could fail to enter config_mode() if there is already a CLI session in this mode.

Sending commands
The most common operation you’re going to perform with Netmiko is sending commands to a device. Let’s look at a few methods to do this.

To simply send a single command to a device, you can use one of three methods:

send_command_expect()
This method is used for long-running commands that may take a while for the device to process (show run on a larger chassis, show tech, etc.). By default, this method waits for the same prompt string to return before completing. Optionally, you can pass what the new prompt string is going to be should it change based on the commands being sent.

send_command_timing()
This method is for short-running commands; it is timing based and does not check the prompt string.

send_command()
This is an older method in Netmiko, which now acts as a wrapper for calling send_command_expect(). Thus, send_command() and send_command_expect() perform the same operation.

Let’s look at a few examples. Here you’re gathering a show run and printing out the first 176 characters for verification:

>>> show_run_output = device.send_command('show run')
>>>
>>> print(show_run_output[:176])

!Command: show running-config
!Running configuration last done at: Wed Oct  5 04:18:12 2022
!Time: Wed Oct  5 04:23:22 2022

version 9.3(3) Bios:version
hostname nxos-spine1
Send a command that changes the prompt string—remember you’re still in configuration mode when you enter device.config_mode()—as follows:

>>> output = device.send_command_expect('end')
Traceback (most recent call last):
  File "<stdin>", line 1, in <module>
  File "/usr/local/lib/python3.8/site-packages/netmiko/base_connection.py",
  line 1582, in send_command_expect
    return self.send_command(*args, **kwargs)
  File "/usr/local/lib/python3.8/site-packages/netmiko/utilities.py", line 500,
   in wrapper_decorator
    return func(self, *args, **kwargs)
  File "/usr/local/lib/python3.8/site-packages/netmiko/base_connection.py",
   line 1535, in send_command
    raise IOError()
OSError: Search pattern never detected in send_command: nxos\-spine1\(config\)\#
>>>
The stack trace shown is expected, as send_command_expect() expects to see the same prompt string by default. Since you are in config mode with the current prompt string of nxos-spine1(config)#, when you type the command end, the new prompt string is going to be nxos-spine1#.

To execute a command that changes the prompt string, you have two options. First, you can use the expect_string parameter that defines the new and expected prompt string:

>>> output = device.send_command_expect('end', expect_string='nxos-spine1#')
>>>
Second, you can use the send_command_timing() method, which is timing based and doesn’t expect a particular prompt string to be found again:

>>> output = device.send_command_timing('end')
>>>
You’ve shown three methods thus far on how to send commands with Netmiko. Let’s look at two more useful ones, as you may want to send several commands at once instead of one at a time.

Netmiko also supports a method called send_config_set() that takes a parameter that must be iterable. We’ll show this using a Python list, but you can also use a Python set:

>>> commands = [
  'interface Ethernet1/1',
  'description configured by netmiko',
  'shutdown'
]
>>>
>>> output = device.send_config_set(config_commands=commands)
>>>
>>> print(output)
nxos-spine1(config)# interface Ethernet1/1
nxos-spine1(config-if)# description configured by netmiko
nxos-spine1(config-if)# shutdown
nxos-spine1(config-if)# end
nxos-spine1#
This method checks whether you’re already in configuration mode. If you aren’t, it goes into config mode, executes the commands, and by default, exits configuration mode. You can verify this by viewing the returned output, as shown in the previous example.

Finally, Netmiko has a method that can execute commands from a file. This allows you to do something like create a Jinja template, render it with variable data, write the data to a file, and then execute those commands from the file with the Netmiko method send_config_from_file(). Building on what we covered in Chapters 6 and 9, let’s see how to perform this workflow in Example 10-38.

Example 10-38. Sending commands from a file with Netmiko
from netmiko import ConnectHandler
from jinja2 import Environment, FileSystemLoader

device = ConnectHandler(
  ...
)

interface_dict = {
    "name": "Ethernet1/2",
    "description": "Server Port",
    "vlan": 10,
    "uplink": False
}

# Create the custom commands, combining the Jinja config.j2 template
# with the data defined in interface_dict
ENV = Environment(loader=FileSystemLoader('.'))
template = ENV.get_template("config.j2")
commands = template.render(interface=interface_dict)

# Store the CLI commands in a local file
filename = 'nxos.conf'
with open(filename, 'w') as config_file:
    config_file.writelines(commands)

# Send CLI commands directly from the file
output = device.send_config_from_file(filename)

# Use show commands to verify that the change succeeded
verification = device.send_command(f'show run interface {interface_dict["name"]}')
print(verification)

device.disconnect()
Everything shown in this example was covered in prior chapters. Note that config.j2 must be created for this to work, and for this example, that the Jinja template is stored in the same directory from where we entered the Python interpreter. The content of the template is from Example 9-3, and is as follows:

interface {{ interface.name }}
 description {{ interface.description }}
 switchport access vlan {{ interface.vlan }}
 switchport mode access
Finally, when you’re done working with Netmiko, you can gracefully disconnect from the device by using the disconnect() method. If we run the script, we will see the verification of the new interface configuration according to the template:

ch10-apis/python_netmiko$ python3 send_commands_from_file.py

!Command: show running-config interface Ethernet1/2
!Running configuration last done at: Wed Oct  5 04:38:54 2022
!Time: Wed Oct  5 04:38:55 2022

version 9.3(3) Bios:version

interface Ethernet1/2
  description Server Port
  switchport access vlan 10
TIP
Context managers in Python help manage setup and teardown operations. For an SSH library such as Netmiko, a context manager seems ideal. Thus, Netmiko provides one, netmiko.ConnectHandler, that will take care of establishing the SSH session at the beginning, and tearing it down when exiting it (so you don’t leave open SSH connections):

with netmiko.ConnectHandler(**device_config) as device:
    device.send_command("show run")
So far, you have shown the benefits of Netmiko, allowing interaction with a traditional CLI interface programmatically. Unfortunately, the unstructured data used in the CLI output is a big drawback in the automation journey. Hopefully, we have some helpers available.

Empowering Netmiko with TextFSM and NTC Templates
TextFSM is an open source project built by Google that converts semiformatted text (the CLI output) to structured data, using templates. So, for each CLI output, you need to provide a specific template that NTC Templates solves.

NTC Templates is an open source project sponsored by Network to Code that provides a large collection of TextFSM templates for a lot of network vendors.

NOTE
You don’t need to install TextFSM or NTC Templates because they are dependencies of Netmiko, so they are already installed.

In Example 10-39, we demonstrate how to use NTC Templates in two steps:

Get raw CLI output from Netmiko and store it as a string.

Use the NTC Templates parser to transform the raw output into structured data.

Example 10-39. Using NTC Templates to get structured data from Netmiko output
>>> from netmiko import ConnectHandler
>>> device = ConnectHandler(
...     host='nxos-spine1',
...     username='admin',
...     password='admin',
...     device_type='cisco_nxos'
... )
>>> show_interfaces_raw = device.send_command('show int brief')
>>> show_interfaces_raw[:150]
'\n-------------------------------------------------------------------------------\n
Port   VRF          Status IP Address                              S'
>>>
>>> from ntc_templates.parse import parse_output
>>> show_interfaces_parsed =  parse_output(
...     platform="cisco_nxos",    1
...     command="show int brief", 2
...     data=show_interfaces_raw, 3
... )
>>> show_interfaces_parsed[0]
{
  'interface': 'mgmt0', 'vrf': '--', 'status': 'up', 'ip': '10.0.0.15',
  'speed': '1000', 'mtu': '1500', 'vlan': '', 'type': '', 'mode': '',
  'reason': '', 'portch': '', 'description': ''
}
1
Indicates the reference platform. Each will have different parsers.

2
Identifies the specific template within a platform because each CLI command may have different data.

3
The raw input data to be parsed.

Luckily, since its 2.0.0 release, Netmiko has the implicit support of NTC Templates, simply using the use_textfsm argument:

>>> show_interfaces_parsed_directly = device.send_command(
...  'show int brief',
...   use_textfsm=True,
... )
>>> show_interfaces_parsed == show_interfaces_parsed_directly
True
This functionality is just combining the two steps presented in Example 10-39.

NOTE
Netmiko is also used as the primary SSH driver for devices within NAPALM, a robust and multivendor network Python library for configuring devices and retrieving data. We cover NAPALM in Chapter 12.

This concludes using Netmiko to automate SSH-based network devices. You’ve now seen how to automate various types of network devices across a range of API types, no matter the device or API type you need to work with.

Summary
This chapter introduced the available types of APIs in the context of networking: HTTP-based APIs (both RESTful and non-RESTful), NETCONF, RESTCONF, and gNMI. After an introduction using command-line tools, we went through how these interfaces can be leveraged with the Go and Python programming languages. All these programmatic interfaces allow a much more efficient automation pattern than the traditional CLI with unstructured data. However, even nowadays, a lot of network devices are still out there where the main interface is still an SSH connection. Because of this, we also covered the Netmiko Python library.

At this point, you may be wondering, “Which network interface should I use?” The answer is, it depends. It depends on the device you are trying to automate—including the interfaces it supports and the capabilities you need to use. Another important factor is the tooling you prefer. Having tools that allow effective and efficient automation is key to success. One way or another, all the data-model interfaces share one big challenge that slows their adoption: the translation between standard or community data models to the device models does not always cover all the features or may not be fully supported. In the end, the goal is to change the data in the device’s data store, and these data models come with inertia. This is the reason text-based configurations are still the main interface for many network devices.

As you continue your journey automating network devices with various types of APIs, remember that there is no magic here—​you need to perform due diligence to understand how to use any given API and to know which options each platform offers over each interface.

In the next chapter, we shift gears and introduce the importance of using source code control to support network automation and programmability.