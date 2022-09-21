
import React, {Component} from 'react'
import SweetAlert from 'react-bootstrap-sweetalert';
import PaymentUI from './PaymentUI';

class ServicePayModal extends Component{ 
    
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
            <SweetAlert show={this.state.show} showConfirm={false} showCancel={false} showCloseButton={true} title={'Place an order'} onCancel={this.onClose}  >
                <PaymentUI purpose='service' user={this.props.user} creator={this.props.creator} onClose={this.onClose} amount={this.props.service.price} service={this.props.service} form={this.props.form}  />
            </SweetAlert>
        )
    }
}


export default ServicePayModal;