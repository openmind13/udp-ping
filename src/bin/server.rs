const ADDR: &str = "0.0.0.0:9000";
const PACKET_SIZE: usize = 508;

fn main() {
    println!("Listen udp on {}", ADDR);

    let udp_listener = std::net::UdpSocket::bind(ADDR).unwrap();
    let mut buf = [0u8; PACKET_SIZE];

    let (amt, src) = udp_listener.recv_from(&mut buf).unwrap();
    println!("{:?} bytes from {}", amt, src);
}
