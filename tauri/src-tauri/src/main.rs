#![cfg_attr(
  all(not(debug_assertions), target_os = "windows"),
  windows_subsystem = "windows"
)]

mod cmd;

use std::process::Command;
use std::thread;

fn start_hq(){
  
  let output = if cfg!(target_os = "windows") {
    Command::new("hq\\hq.exe")
            .args(&["-c", "config.yml"])
            .output()
            .expect("failed to execute process")
  } else {
    Command::new("sh")
            .arg("-c")
            .arg("echo hello")
            .output()
            .expect("failed to execute process")
  };

  println!("status: {}", output.status);
  println!("stdout: {}", String::from_utf8_lossy(&output.stdout));
  println!("stderr: {}", String::from_utf8_lossy(&output.stderr));

}

fn start_hq1(){

  const NTHREADS: u32 = 10;
  let mut children = vec![];

  for i in 0..NTHREADS {
      // Spin up another thread
      children.push(thread::spawn(move || {
          println!("this is thread number {}", i);
      }));
  }

  for child in children {
      // Wait for the thread to finish. Returns a result.
      let _ = child.join();
  }

}

fn main() {
  // thread::spawn(move || { start_hq(); });
  if cfg!(target_os = "windows") {
    Command::new("hq\\hq.exe").args(&["-c", "config.yml"]).spawn();
  }

  println!("{}", "main.rs");
  tauri::AppBuilder::new()
    .invoke_handler(|_webview, arg| {
      use cmd::Cmd::*;
      match serde_json::from_str(arg) {
        Err(e) => {
          Err(e.to_string())
        }
        Ok(command) => {
          match command {
            // definitions for your custom commands from Cmd here
            MyCustomCommand { argument } => {
              //  your command code
              println!("{}", argument);
            }
          }
          Ok(())
        }
      }
    })
    .build()
    .run();
}
