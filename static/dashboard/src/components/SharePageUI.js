
import React, { Component } from 'react'
import v1 from '../api/v1'


function fallbackCopyTextToClipboard(text) {
  var textArea = document.createElement("textarea");
  textArea.value = text;
  
  // Avoid scrolling to bottom
  textArea.style.top = "0";
  textArea.style.left = "0";
  textArea.style.position = "fixed";

  document.body.appendChild(textArea);
  textArea.focus();
  textArea.select();

  try {
    var successful = document.execCommand('copy');
    var msg = successful ? 'successful' : 'unsuccessful';
    console.log('Fallback: Copying text command was ' + msg);
  } catch (err) {
    console.error('Fallback: Oops, unable to copy', err);
  }

  document.body.removeChild(textArea);
}
function copyTextToClipboard(text) {
  if (!navigator.clipboard) {
    fallbackCopyTextToClipboard(text);
    return;
  }
  navigator.clipboard.writeText(text).then(function() {
    console.log('Async: Copying to clipboard was successful!');
  }, function(err) {
    console.error('Async: Could not copy text: ', err);
  });
}


async function sharePage(evt) {
        try {
            copyTextToClipboard(window.location.href)
            
        } finally {
            evt.preventDefault()
        }
}
    
class SharePageUI extends Component{
  constructor(props) {
    super(props)
    v1.page.event('Share Modal', "Open", '')

    this.state = {
      copyBtnText: 'Copy',
      copyBtnClasses: '',
    }

    this.share = evt => {
      
      sharePage(evt)
      this.setState({copyBtnText: 'Copied!', copyBtnClasses: "bg-danger"})
      setTimeout(_ => {
        this.setState({copyBtnText: 'Copy', copyBtnClasses: ""})
      }, 2500)
    }
  }
  
  shareLink(target) {
    let url = window.location.href
    switch (target) {
      case "facebook":
        return "https://www.facebook.com/sharer/sharer.php?u="+ url
      case "twitter":
        return "https://twitter.com/intent/tweet?url=" + url
      case "whatsapp":
        return "https://wa.me/?text=" + url
      default:
        return url
    }
  }
   
  render() {

    return (
      <div class="modal-body">
        <p>Share this link via</p>
        <div class="d-flex align-items-center icons">
          <a href={this.shareLink("facebook")} target="_blank" rel="noopener" class="fs-5 d-flex align-items-center justify-content-center">
            <span class="fab fa-facebook-f"></span>
          </a>
          <a href={this.shareLink("twitter")} target="_blank" rel="noopener" class="fs-5 d-flex align-items-center justify-content-center">
            <span class="fab fa-twitter"></span>
          </a>
          <a href={this.shareLink("whatsapp")} target="_blank" rel="noopener" class="fs-5 d-flex align-items-center justify-content-center">
            <span class="fab fa-whatsapp"></span>
          </a>
        </div>
        <p>Or copy link</p>
        <div class="field d-flex align-items-center justify-content-between">
          <span class="fas fa-link text-center"></span>
          <input type="text" value={window.location.href} />
          <button class={this.state.copyBtnClasses} onClick={this.share}>{this.state.copyBtnText}</button>
        </div>
      </div>
    )
  }
}
export default SharePageUI;