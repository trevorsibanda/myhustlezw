
import React, {Component} from 'react'
import SweetAlert from 'react-bootstrap-sweetalert';
import SharePageUI from './SharePageUI';

class SharePageModal extends Component{ 
    
    constructor(props) {
        super(props);
        this.state = {
            show: true,
        }
        this.onClose = this.onClose.bind(this)
    }

    

    onClose() {
        this.setState({
            show: false
        })
        return this.props.onClose ? this.props.onClose() : null
    }

    render() {
        return (
            <SweetAlert show={this.state.show} showConfirm={false} showCancel={false} showCloseButton={true} title={'Share this page'} onCancel={this.onClose}  >
                <SharePageUI />
            </SweetAlert>
        )
    }
}


export default SharePageModal;