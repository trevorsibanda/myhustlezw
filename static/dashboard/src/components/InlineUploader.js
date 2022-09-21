import React, { Component } from "react"

import Uppy from '@uppy/core'
import { Dashboard, FileInput } from '@uppy/react'
import XHRUpload from '@uppy/xhr-upload'

import '@uppy/core/dist/style.css'
import '@uppy/dashboard/dist/style.css'

import v1 from "../api/v1";
import ImageEditor from "@uppy/image-editor"


class InlineUpload extends Component {

    constructor(props){

        super(props)
        this.state = {
            uploadOpen: true,
        }

        this.map = {}

        this.upload = new Uppy({
            allowMultipleUploads: this.props.maxNumberOfFile > 1,
            restrictions: {
                maxNumberOfFiles: this.props.maxNumberOfFiles ? this.props.maxNumberOfFiles : 1,
                minFileSize: (1024 * 1),
                maxFileSize: (1024 * 1024 * 1024 * 2),
                allowedFileTypes: this.props.allowedTypes,
            }
        })
            .use(XHRUpload, { endpoint: v1.config.uploadMediaEndpoint(this.props.type, this.props.purpose) })
            .use(ImageEditor, {
                id: 'ImageEditor',
                quality: 0.8,
                cropperOptions: {
                viewMode: 1,
                background: false,
                autoCropArea: 1,
                responsive: true,
                },
            })
            
        
        this.upload.on('upload-success', (file, response) => {
            if (response.body._id) {
                this.map[file.id] = response.body
                return this.props.onUploaded ? this.props.onUploaded(response.body) : console.log(response.body)
            } else {
                alert('Upload failed with error\n:'+ response.body.error)
            }
            this.upload.reset()
        })

        this.upload.on('file-removed', (file, reason) => {
            this.upload.removeFile(file.id)
            v1.files.delete(this.map[file.id]).then(_ => {
                console.log('Deleted file ', file.id)
            }).catch(err => {
                console.error('Failed to delete file with error: ', file.id,  err)
            })
        })

        if (!this.props.maxNumberOfFile || this.props.maxNumberOfFile <= 1) {
            console.log('set it ups')
            this.upload.on('file-added', file => {
                if(!file.type.startsWith('image')){
                    return
                }
                console.log('file-added')
            const dashboard = this.upload.getPlugin('react:Dashboard')
                console.log(dashboard)
                window.ufile = file
            //dashboard.toggleFileCard(true, file.id)
                try {
                    setTimeout(_ => dashboard.openFileEditor(file), 250)
                } catch (e) {
                    console.log('exception')
                    console.log(e)
                }
                console.log('opened file editor')
            })
        }

        window.uppy = this.upload

        this.handleOpenF = () => {
            this.setState({ uploadOpen: true })
        }

        this.handleCloseF = () => {
            this.setState({ uploadOpen: false })
            this.upload.reset()
        }

    }

    componentWillUnmount(){
        this.upload.close()
    }

    render() {
        return this.props.noDashboard ? <FileInput uppy={this.upload} /> : (
            <Dashboard
                uppy={this.upload}
                showRemoveButtonAfterComplete={true}
                doneButtonHandler={this.props.doneBtnHandler ? this.props.doneBtnHandler : () => { }}
                open={this.state.uploadOpen}
                onRequestClose={this.handleCloseF}
                plugins={['ImageEditor']}
            />
        )
    }
}

export default InlineUpload;