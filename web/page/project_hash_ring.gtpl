{{define "page"}}
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <title>Hash Ring - Hypermind</title>
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta name="description" content="The homepage of Hypermind.">
    <meta name="author" content="hyper-carrot">

    {{template "header-import"}}
    {{template "js-import"}}

</head>
<body>

{{template "top-navbar" .}}

<div class="container-fluid">
    <div class="row-fluid">
        <div class="span2">
            {{template "projects-navbar" .}}
        </div>
        <div class="span10">
            <div class="hero-unit">
                <p>
                    It from the internal training by me speech in my company.
                </p>
                <p>
                    Hash ring is a kind of consistency hash realization. 
                    Now, consider this situation: A software system was added a data cache layer in order to relieve the database of pressure as a result of a mass of query. In data cache layer, they use multi nodes cluster operation mode. The problem: How data in these nodes distribution? How about the performance, expandability and complexity?
                </p>
                <p> 
                    <h3>The Classical Solution:</h3>
                    There are 3~5 cache nodes.
                    <br>
                    They distribute data to the cache nodes using "hash (keyword) mod n" strategy.
                    <br>
                    <img src="/img/chash-hash+mod.png">
                    <br>
                    The disadvantage:
                    <br>
                    <ul>
                        <li>If the cache node is down, all the operations for the data assigned to the cache node will be unavailable.</li>
                        <li>The horizontal scalability for the cluster is difficult.</li>
                        <li>The uneven distribution data  is not easy to be adjusted.</li>
                    </ul>
                </p>
                <p> 
                    <h3>Hash ring - the standard configuration for consistency hash:</h3>
                    What is the hash ring?
                    <br>
                    <ul>
                        <li>A continuous, covering positive integer range of hash range, and the first integer meeting the last.</li>
                        <li>The cache node hash value (the URI hash value) and the data of the hash value (the keyword hash value) are stored in the ring.</li>
                    </ul>
                    The blank ring:
                    <br>
                    <img src="/img/chash-hash_ring-intro1.png">
                    <br>
                    Let's see the situation after adding cache node:
                    <br>
                    <img src="/img/chash-hash_ring-intro2.png">
                    <br>
                    Now, I add some data nodes:
                    <br>
                    <img src="/img/chash-hash_ring-intro3.png">
                    <br>
                    As shown, 
                    <ol>
                        <li>The data A stored in cache node A.</li>
                        <li>The data B stored in cache node B.</li>
                        <li>The data C stored in cache node D.</li>
                        <li>The data D stored in cache node D.</li>
                    </ol> 
                    The steps for find the target cache node:
                    <ol>
                        <li>Select a keyword in data.</li>
                        <li>Using the keyword 'sha1' value as the data hash.</li>
                        <li>Query the hash ring, and find out the cache node which its hash greater than & neighboring the data hash.</li>
                        <li>Store the data in the cache node.</li>
                    </ol>
                    (The location method for data getting is alike.)
                    <br>
                    Now, the cache B is down!
                    <br>
                    <img src="/img/chash-hash_ring-intro4.png">
                    <br>
                    The data B now will be stored in the cache node C, not cache node B. The hash ring can timed check the availability of every cache node , and remove the cache node when it is unavailable.
                    <br>
                    Next, We add a cache node named 'Cache E'.
                    <br>
                    <img src="/img/chash-hash_ring-intro5.png">
                    <br>
                    Then, The data D now will be stored in the cache node E, not cache node D. Because the hash of 'Cache E' is greater than & neighboring the hash of 'Data D'.
                    <br>
                    With this, some problems are resolved:
                    <ul>
                        <li>When a cache node is down machine, the corresponding data will be automatic transfered to other available nodes.</li>
                        <li>When need to add a cache node, the part of the data on a existing cache node in hash ring will be  automatic transfered to add new cache nodes.</li>
                    </ul>
                    However, there are some problems yet:
                    <ul>
                        <li>When a cache node is down, the data will be transfered automatically. But the migration pressure will all point to a certain cache node. This operation may cause this cache node data volume and pressure spurt, even is down.</li>
                        <li>When need to add a cache node, only will share a certain existing cache node data volume and pressure. The advantage of this seemingly is not obvious.</li>
                        <li>The data distribution is not uniform, because the cache nodes are few.</li>
                    </ul>
                </p>
                <p> 
                    <h3>Hash ring with virtual node:</h3>
                    Q: What is virtual node?
                    <br>
                    A: A cache server is not only the corresponding a cache node in hash ring, but more virtual nodes. They are uniform distribution in the hash ring.
                    <br>
                    Q: How to do it?
                    <br>
                    A: The shadow algorithm: 
                    <ol>
                        <li>A keyword expansion into more than one keyword</li>
                        <li>KETAMA algorithm:</li>
                        <img src="/img/ketama_python.png">
                    </ol>
                    <br>
                    Q: The virtual node means what?
                    <br>
                    A: This means that:
                    <ol>
                        <li>When a cache server is down, the data on its virtual node will were migrated to many other cache servers.</li>
                        <li>When need to add a cache server, its 'shadows' will share data volume and pressure from different existing cache nodes.</li>
                    </ol>
                    <img src="/img/chash-hash_ring-intro6.png">
                    <br>
                    In conclusion, the hash ring with virtual node can more uniform distribute data volume and pressure on the cluster.
                    <br>
                </p>
                <p>
                    I had implemented hash ring model written by Go, Python and Java. And, I compared their performance. (see picture below)
                    <br>
                    <img src="/img/chash_benchmark2.png">
                    <br>
                    The source code of hash ring implemented by me is here: <br>
                    Go Edition: <a href="https://github.com/hyper-carrot/chash4go" target="_blank">chash4go</a><br>
                    Python Edition: <a href="https://github.com/hyper-carrot/chash4py" target="_blank">chash4py</a><br>
                    Java Edition: <a href="https://github.com/hyper-carrot/chash4j" target="_blank">chash4j</a><br>
                    All of these editions can be used in the production environment.
                </p>
                <p>
                    That's all. Welcome to exchange!
                </p>
            </div>
        </div>
    </div>
</div>

</body>
</html>
{{end}}