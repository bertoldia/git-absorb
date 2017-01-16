extern crate git2;

use std::env;
use std::process::Command;
use git2::Repository;
use git2::string_array::StringArray;

fn print_remotes(remotes: StringArray) {
    for remote in remotes.iter() {
        match remote {
            Some(rem) => println!("{}", rem),
            None => println!("Shitbox"),
        }
    }
}

fn main() {
    let target_commit = env::args().nth(1).expect("Missing target commit.");

    let repo = match Repository::open(".") {
        Ok(repo) => repo,
        Err(e) => panic!("failed to open: {}", e),
    };

    println!("{:?}", repo.state());

    match repo.remotes() {
        Ok(rem) => print_remotes(rem),
        Err(e) => panic!("Oh shit! {}", e),
    };

    match repo.revparse_single(&target_commit) {
        Ok(spec) => println!("{:?}", &spec.id()),
        Err(e) => panic!("Invalid revspec {}: {}!", target_commit, e),
    };

    let output = Command::new("git")
        .arg("rev-parse")
        .arg(&target_commit)
        .output()
        .expect("failed to execute process");

    println!("status: {}", &output.status);
    println!("stdout: {}", String::from_utf8_lossy(&output.stdout));
    println!("stderr: {}", String::from_utf8_lossy(&output.stderr));
}
