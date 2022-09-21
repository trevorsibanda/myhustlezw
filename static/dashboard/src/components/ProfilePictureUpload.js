import {Component} from "react"

import Uppy from '@uppy/core'
import { DashboardModal } from '@uppy/react'
import XHRUpload from '@uppy/xhr-upload'

import '@uppy/core/dist/style.css'
import '@uppy/dashboard/dist/style.css'
import '@uppy/webcam/dist/style.css'


import v1 from "../api/v1"

class ProfilePictureUpload extends Component {
    constructor(props){
        super(props)

        this.state = {
            modalOpen: false,
            profpicurl: v1.assets.profPicURL(this.props.user)
        }

        this.image = new Uppy({
            allowMultipleUploads: false,
            restrictions: {
                maxNumberOfFiles: 1,
                minFileSize: (1024 * 100),
                maxFileSize: (1024 * 1024 * 10),
                allowedFileTypes: ['image/*']
            }
        }).use(XHRUpload, { endpoint: v1.config.uploadMediaEndpoint('image') })
        

        this.image.on('upload-success', (file, response) => {
            if (response.body._id){
                v1.user.setAvatar(response.body._id).then(resp => {
                    if (this.props.onUploaded) {
                        this.props.onUploaded(file, response)
                    }
                    this.handleClose()
                })
                this.setState({
                    profpicurl: this.state.profpicurl + '?reload'
                })
                
            } else {
                alert('Failed to set new profile picture. Please upload a new image')
            }
        })


        this.handleOpen = this.handleOpen.bind(this)
        this.handleClose = this.handleClose.bind(this)
    }

    componentWillUnmount() {
        this.image.close()
    }

    handleOpen() {
        this.setState({ modalOpen: true })
    }

    handleClose() {
        //check if upload in progress
        this.setState({ modalOpen: false })
        this.image.reset()
    }
    
    render(){
        return (
            <>
                <DashboardModal
                    uppy={this.image}
                    closeModalOnClickOutside
                    open={this.state.modalOpen}
                    onRequestClose={this.handleClose}
                />
                <a href="#" onClick={this.handleOpen} >
                    <img title="Click to upload new profile picture"
                        className="avatar avatar-xxl avatar-bordered"
                        src={this.state.profpicurl} alt=""
                        style={{"width":"15vh", "height": "15vh"}}
                    />
                </a>
                
            </>
        )
    }
}

export default ProfilePictureUpload;