# world-of-wisdom

I choose hashcash algorithm, because of:

Proven Security: Hashcash has a proven track record of security and has been widely used in various applications, including email spam prevention and cryptocurrency mining. It relies on the computational effort required to find a hash collision, making it resistant to brute-force attacks.

Efficiency: Hashcash offers a good balance between computational effort and the time it takes to verify the PoW. It is computationally intensive to find the correct hash, but verification is relatively quick and straightforward, making it an efficient choice.

idely Adopted: Hashcash is a well-known and widely adopted PoW system, which means that there are existing libraries and tools available for its implementation in various programming languages, including Go.

Customization: You can customize the difficulty level of the PoW challenge by adjusting the number of leading zeros required in the hash, allowing you to control the level of computational effort required for validation.

Scalability: Hashcash can be adjusted to suit different use cases and scalability requirements. Whether you're implementing it for anti-DDoS protection or other purposes, you can adapt the PoW parameters to your specific needs.

Resistance to Parallelism: Hashcash's reliance on finding hash collisions is inherently resistant to parallelization, which can deter attackers from easily scaling their computational power to break the system.

Energy Efficiency: Compared to some other PoW algorithms, Hashcash is relatively energy-efficient, making it a suitable choice for applications that aim to minimize environmental impact.

## how to run 

``make run``

It will build docker images and run server and client. 