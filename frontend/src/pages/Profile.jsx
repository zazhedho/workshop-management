import React, { useState } from 'react'
import { Card, Form, Button, Alert, Row, Col } from 'react-bootstrap'
import { useAuth } from '../contexts/AuthContext'

const Profile = () => {
  const { user, updateProfile } = useAuth()
  const [formData, setFormData] = useState({
    name: user?.name || '',
    email: user?.email || '',
    phone: user?.phone || '',
    password: '',
    confirmPassword: ''
  })
  const [error, setError] = useState('')
  const [success, setSuccess] = useState('')
  const [loading, setLoading] = useState(false)

  const handleChange = (e) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value
    })
  }

  const handleSubmit = async (e) => {
    e.preventDefault()
    setError('')
    setSuccess('')

    if (formData.password && formData.password !== formData.confirmPassword) {
      setError('Passwords do not match')
      return
    }

    setLoading(true)

    const updateData = {
      name: formData.name,
      email: formData.email,
      phone: formData.phone
    }

    if (formData.password) {
      updateData.password = formData.password
    }

    const result = await updateProfile(updateData)
    
    if (result.success) {
      setSuccess('Profile updated successfully')
      setFormData(prev => ({
        ...prev,
        password: '',
        confirmPassword: ''
      }))
    } else {
      const errorPayload = result.error
      if (errorPayload && errorPayload.message) {
        setError(errorPayload.message)
      } else if (errorPayload) {
        setError(String(errorPayload))
      } else {
        setError('Profile update failed')
      }
    }
    
    setLoading(false)
  }

  return (
    <div>
      <div className="d-flex justify-content-between align-items-center mb-4">
        <h2>Profile Settings</h2>
      </div>

      <Row>
        <Col md={8} lg={6}>
          <Card>
            <Card.Header>
              <h5 className="mb-0">Update Profile Information</h5>
            </Card.Header>
            <Card.Body>
              {error && <Alert variant="danger">{error}</Alert>}
              {success && <Alert variant="success">{success}</Alert>}

              <Form onSubmit={handleSubmit}>
                <Form.Group className="mb-3">
                  <Form.Label>Full Name</Form.Label>
                  <Form.Control
                    type="text"
                    name="name"
                    value={formData.name}
                    onChange={handleChange}
                    required
                    placeholder="Enter your full name"
                  />
                </Form.Group>

                <Form.Group className="mb-3">
                  <Form.Label>Email</Form.Label>
                  <Form.Control
                    type="email"
                    name="email"
                    value={formData.email}
                    onChange={handleChange}
                    required
                    placeholder="Enter your email"
                  />
                </Form.Group>

                <Form.Group className="mb-3">
                  <Form.Label>Phone</Form.Label>
                  <Form.Control
                    type="tel"
                    name="phone"
                    value={formData.phone}
                    onChange={handleChange}
                    required
                    placeholder="Enter your phone number"
                  />
                </Form.Group>

                <hr />

                <h6 className="mb-3">Change Password (Optional)</h6>

                <Form.Group className="mb-3">
                  <Form.Label>New Password</Form.Label>
                  <Form.Control
                    type="password"
                    name="password"
                    value={formData.password}
                    onChange={handleChange}
                    placeholder="Enter new password (leave blank to keep current)"
                  />
                </Form.Group>

                <Form.Group className="mb-4">
                  <Form.Label>Confirm New Password</Form.Label>
                  <Form.Control
                    type="password"
                    name="confirmPassword"
                    value={formData.confirmPassword}
                    onChange={handleChange}
                    placeholder="Confirm new password"
                  />
                </Form.Group>

                <Button
                  type="submit"
                  variant="primary"
                  disabled={loading}
                  className="w-100"
                >
                  {loading ? (
                    <>
                      <i className="fas fa-spinner fa-spin me-2"></i>
                      Updating...
                    </>
                  ) : (
                    'Update Profile'
                  )}
                </Button>
              </Form>
            </Card.Body>
          </Card>
        </Col>

        <Col md={4} lg={6}>
          <Card>
            <Card.Header>
              <h5 className="mb-0">Account Information</h5>
            </Card.Header>
            <Card.Body>
              <div className="text-center mb-4">
                <i className="fas fa-user-circle fa-5x text-muted"></i>
              </div>
              
              <div className="mb-3">
                <strong>Name:</strong>
                <p className="mb-0">{user?.name}</p>
              </div>
              
              <div className="mb-3">
                <strong>Email:</strong>
                <p className="mb-0">{user?.email}</p>
              </div>
              
              <div className="mb-3">
                <strong>Phone:</strong>
                <p className="mb-0">{user?.phone}</p>
              </div>
              
              <div className="mb-3">
                <strong>Role:</strong>
                <p className="mb-0">
                  <span className="badge bg-primary">{user?.role}</span>
                </p>
              </div>
              
              <div className="mb-3">
                <strong>Member Since:</strong>
                <p className="mb-0">
                  {user?.created_at && new Date(user.created_at).toLocaleDateString('en-US', {
                    year: 'numeric',
                    month: 'long',
                    day: 'numeric'
                  })}
                </p>
              </div>
            </Card.Body>
          </Card>
        </Col>
      </Row>
    </div>
  )
}

export default Profile