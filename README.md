<center>

![logo](_img/logo.png)

</center>

<h4 align="center"> 📡 ScopeDetective - Tool for Monitoring </h4>
<p align="center">
  <a href="#features">Features</a> •
  <a href="#installation">Installation</a> •
  <a href="#usage">Usage</a> •
  <a href="#future">Future</a> •
  <a href="#contribute">Contribute</a> •
  <a href="#contact">Contact me</a>
</p>

---

ScopeDetective is a rapid tool for monitoring scope changes on HackerOne. It provides a simple and efficient way to track modifications to the scope of your bug bounty program. By using ScopeDetective, you can stay up-to-date with any changes and ensure that you are aware of all the assets in scope.


## Features
* Rapidly monitors scope changes on the HackerOne platform
* Sends notifications to a specified Discord webhook
* Allows users to set the delay time between each monitoring check


## Installation

[Download](https://github.com/NImaism/ScopeDetective/releases/latest) a prebuilt binary from [releases page](https://github.com/NImaism/ScopeDetective/releases/latest), unpack and run!

_or_

If you have recent go compiler installed: `go install github.com/NImaism/ScopeDetective@latest` (the same command works for updating)

_or_

`git clone https://github.com/NImaism/ScopeDetective ; cd ScopeDetective ; go get ; go build`
<br/>

Note: ScopeDetective depends on Go **1.20** or **greater**.


## Usage
To use ScopeDetective, follow these steps:
1. Open your command-line interface (CLI).
2. Run the following command: `ScopeDetective -webhook <DiscordWebhookUrl> -delay <delayTime>`
3. Replace `<DiscordWebhookUrl>` with the actual URL of your Discord webhook.
4. Replace `<delayTime>` with the desired delay time between each monitoring check. Specify the delay in minutes.

Example: `ScopeDetective -webhook https://discord.com/webhook -delay 5`
By running this command, ScopeDetective will start monitoring scope changes and send notifications to your specified Discord webhook with the specified delay between each check.


## Future

We have several exciting plans for the future development of ScopeDetective, including:

* Adding support for monitoring other bug bounty platforms
* Enhancing the user interface for a more user-friendly experience
* Implementing additional notification channels, such as Slack and email
  Stay tuned for updates and new features in upcoming releases!


## Contribute

We welcome contributions from the open-source community to enhance ScopeDetective. If you'd like to contribute, please follow these steps:

1. Create an issue: Before making any changes, please create a new issue to discuss your proposed changes, bug fixes, or new features. This allows for better coordination and feedback from the project maintainers.

2. Fork the ScopeDetective repository and clone it to your local machine.

3. Make your desired changes, improvements, or bug fixes.

4. Submit a pull request: Once you've made your changes, submit a pull request. Please provide a clear description of your changes, reference the relevant issue, and explain the motivation behind your contribution.

We appreciate your contributions and look forward to collaborating with you to make ScopeDetective even better!


## Contact

If you have any questions or need further assistance, please don't hesitate to reach out. Happy monitoring with [Me](mailto:nima.gholamyy@gmail.com)!


## License

ScopeDetective is released under MIT license. See [LICENSE](https://github.com/NImaism/ScopeDetective/blob/master/LICENSE).


