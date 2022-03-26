const PACKET_SIZE: usize = 508;

fn main() {
    println!("udp client");

    let socket = std::net::UdpSocket::bind("0.0.0.0:9001").unwrap();
    let shared_socket = std::sync::Arc::new(socket);

    let receiver = shared_socket.clone();
    let sender = shared_socket.clone();

    let packet = "test".as_bytes();
    let remote_addr = "0.0.0.0:9000";
    std::thread::spawn(move || loop {
        let n = sender.send_to(&packet, &remote_addr).unwrap();
        println!("Write {n} bytes to {remote_addr}");
        std::thread::sleep(std::time::Duration::from_secs(1))
    });

    let mut recv_buf = [0u8; PACKET_SIZE];
    loop {
        let (n, addr) = receiver.recv_from(&mut recv_buf).unwrap();
        println!("Read {n} bytes from {addr}");
    }
}
