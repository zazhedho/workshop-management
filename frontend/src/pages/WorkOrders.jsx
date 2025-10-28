import React, { useState, useEffect } from 'react'
import { Card, Button, Table, Badge, Modal, Form, Alert, Spinner, Row, Col } from 'react-bootstrap'
import api from '../services/api'

const WorkOrders = () => {
  const [workOrders, setWorkOrders] = useState([])
  const [bookings, setBookings] = useState([])
  const [vehicles, setVehicles] = useState([])
  const [services, setServices] = useState([])
  const [mechanics, setMechanics] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [success, setSuccess] = useState('')

  const [showCreateModal, setShowCreateModal] = useState(false)
  const [showDetailModal, setShowDetailModal] = useState(false)
  const [showAssignModal, setShowAssignModal] = useState(false)
  const [selectedWorkOrder, setSelectedWorkOrder] = useState(null)

  const [createFormData, setCreateFormData] = useState({
    booking_id: '',
    notes: ''
  })

  const [assignMechanicData, setAssignMechanicData] = useState({
    mechanic_id: ''
  })

  useEffect(() => {
    fetchWorkOrders()
    fetchBookings()
    fetchVehicles()
    fetchServices()
    fetchMechanics()
  }, [])

  const fetchWorkOrders = async () => {
    try {
      setLoading(true)
      const response = await api.get('/work-order')
      setWorkOrders(response.data.data || [])
      setError('')
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to fetch work orders')
    } finally {
      setLoading(false)
    }
  }

  const fetchBookings = async () => {
    try {
      const response = await api.get('/booking')
      setBookings(response.data.data || [])
    } catch (err) {
      console.error('Failed to fetch bookings:', err)
    }
  }

  const fetchVehicles = async () => {
    try {
      const response = await api.get('/vehicle')
      setVehicles(response.data.data || [])
    } catch (err) {
      console.error('Failed to fetch vehicles:', err)
    }
  }

  const fetchServices = async () => {
    try {
      const response = await api.get('/service')
      setServices(response.data.data || [])
    } catch (err) {
      console.error('Failed to fetch services:', err)
    }
  }

  const fetchMechanics = async () => {
    try {
      const response = await api.get('/user')
      const allUsers = response.data.data || []
      setMechanics(allUsers.filter(user => user.role === 'mechanic'))
    } catch (err) {
      console.error('Failed to fetch mechanics:', err)
    }
  }

  const handleCreateWorkOrder = async (e) => {
    e.preventDefault()
    setError('')
    setSuccess('')

    try {
      await api.post('/work-order', createFormData)
      setSuccess('Work order created successfully')
      setShowCreateModal(false)
      setCreateFormData({ booking_id: '', notes: '' })
      fetchWorkOrders()
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to create work order')
    }
  }

  const handleViewDetails = async (workOrder) => {
    try {
      const response = await api.get(`/work-order/${workOrder.id}`)
      setSelectedWorkOrder(response.data.data)
      setShowDetailModal(true)
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to fetch work order details')
    }
  }

  const handleAssignMechanic = async (e) => {
    e.preventDefault()
    setError('')
    setSuccess('')

    try {
      await api.put(`/work-order/${selectedWorkOrder.id}/assign`, assignMechanicData)
      setSuccess('Mechanic assigned successfully')
      setShowAssignModal(false)
      setAssignMechanicData({ mechanic_id: '' })
      fetchWorkOrders()
      if (showDetailModal) {
        handleViewDetails(selectedWorkOrder)
      }
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to assign mechanic')
    }
  }

  const handleUpdateStatus = async (workOrderId, newStatus) => {
    setError('')
    setSuccess('')

    try {
      await api.put(`/work-order/${workOrderId}/status`, { status: newStatus })
      setSuccess(`Work order status updated to ${newStatus}`)
      fetchWorkOrders()
      if (showDetailModal && selectedWorkOrder?.id === workOrderId) {
        handleViewDetails({ id: workOrderId })
      }
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to update status')
    }
  }

  const getStatusBadge = (status) => {
    const variants = {
      pending: 'warning',
      in_progress: 'info',
      completed: 'success',
      cancelled: 'danger'
    }
    return <Badge bg={variants[status] || 'secondary'}>{status?.replace('_', ' ').toUpperCase()}</Badge>
  }

  const getBookingInfo = (bookingId) => {
    const booking = bookings.find(b => b.id === bookingId)
    return booking ? `Booking #${booking.id.slice(0, 8)}` : 'N/A'
  }

  const getVehicleInfo = (vehicleId) => {
    const vehicle = vehicles.find(v => v.id === vehicleId)
    return vehicle ? `${vehicle.brand} ${vehicle.model} (${vehicle.license_plate})` : 'N/A'
  }

  const getMechanicName = (mechanicId) => {
    const mechanic = mechanics.find(m => m.id === mechanicId)
    return mechanic ? mechanic.name : 'Unassigned'
  }

  const calculateTotal = (workOrder) => {
    const servicesTotal = workOrder?.services?.reduce((sum, svc) => sum + (svc.price * svc.quantity), 0) || 0
    const partsTotal = workOrder?.parts?.reduce((sum, part) => sum + (part.price * part.quantity), 0) || 0
    return servicesTotal + partsTotal
  }

  if (loading) {
    return (
      <div className="text-center mt-5">
        <Spinner animation="border" role="status">
          <span className="visually-hidden">Loading...</span>
        </Spinner>
      </div>
    )
  }

  return (
    <div>
      <div className="d-flex justify-content-between align-items-center mb-4">
        <h2>
          <i className="fas fa-clipboard-list me-2"></i>
          Work Orders
        </h2>
        <Button variant="primary" onClick={() => setShowCreateModal(true)}>
          <i className="fas fa-plus me-2"></i>
          New Work Order
        </Button>
      </div>

      {error && <Alert variant="danger" dismissible onClose={() => setError('')}>{error}</Alert>}
      {success && <Alert variant="success" dismissible onClose={() => setSuccess('')}>{success}</Alert>}

      <Card>
        <Card.Body>
          <Table responsive hover>
            <thead>
              <tr>
                <th>ID</th>
                <th>Booking</th>
                <th>Vehicle</th>
                <th>Mechanic</th>
                <th>Status</th>
                <th>Created</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              {workOrders.length === 0 ? (
                <tr>
                  <td colSpan="7" className="text-center text-muted py-4">
                    No work orders found
                  </td>
                </tr>
              ) : (
                workOrders.map((workOrder) => (
                  <tr key={workOrder.id}>
                    <td>
                      <small className="font-monospace">{workOrder.id.slice(0, 8)}</small>
                    </td>
                    <td>{getBookingInfo(workOrder.booking_id)}</td>
                    <td>{getVehicleInfo(workOrder.vehicle_id)}</td>
                    <td>{getMechanicName(workOrder.mechanic_id)}</td>
                    <td>{getStatusBadge(workOrder.status)}</td>
                    <td>{new Date(workOrder.created_at).toLocaleDateString()}</td>
                    <td>
                      <Button
                        variant="outline-primary"
                        size="sm"
                        className="me-2"
                        onClick={() => handleViewDetails(workOrder)}
                      >
                        <i className="fas fa-eye"></i>
                      </Button>
                      <Button
                        variant="outline-secondary"
                        size="sm"
                        className="me-2"
                        onClick={() => {
                          setSelectedWorkOrder(workOrder)
                          setShowAssignModal(true)
                        }}
                      >
                        <i className="fas fa-user-plus"></i>
                      </Button>
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </Table>
        </Card.Body>
      </Card>

      <Modal show={showCreateModal} onHide={() => setShowCreateModal(false)} size="lg">
        <Modal.Header closeButton>
          <Modal.Title>Create Work Order</Modal.Title>
        </Modal.Header>
        <Form onSubmit={handleCreateWorkOrder}>
          <Modal.Body>
            <Form.Group className="mb-3">
              <Form.Label>Booking</Form.Label>
              <Form.Select
                value={createFormData.booking_id}
                onChange={(e) => setCreateFormData({ ...createFormData, booking_id: e.target.value })}
                required
              >
                <option value="">Select a booking</option>
                {bookings.filter(b => b.status === 'confirmed').map(booking => (
                  <option key={booking.id} value={booking.id}>
                    Booking #{booking.id.slice(0, 8)} - {getVehicleInfo(booking.vehicle_id)} - {new Date(booking.booking_date).toLocaleDateString()}
                  </option>
                ))}
              </Form.Select>
            </Form.Group>

            <Form.Group className="mb-3">
              <Form.Label>Notes</Form.Label>
              <Form.Control
                as="textarea"
                rows={3}
                placeholder="Enter any additional notes"
                value={createFormData.notes}
                onChange={(e) => setCreateFormData({ ...createFormData, notes: e.target.value })}
              />
            </Form.Group>
          </Modal.Body>
          <Modal.Footer>
            <Button variant="secondary" onClick={() => setShowCreateModal(false)}>
              Cancel
            </Button>
            <Button variant="primary" type="submit">
              Create Work Order
            </Button>
          </Modal.Footer>
        </Form>
      </Modal>

      <Modal show={showDetailModal} onHide={() => setShowDetailModal(false)} size="lg">
        <Modal.Header closeButton>
          <Modal.Title>Work Order Details</Modal.Title>
        </Modal.Header>
        <Modal.Body>
          {selectedWorkOrder && (
            <>
              <Row className="mb-3">
                <Col md={6}>
                  <strong>Work Order ID:</strong>
                  <p className="font-monospace">{selectedWorkOrder.id}</p>
                </Col>
                <Col md={6}>
                  <strong>Status:</strong>
                  <p>{getStatusBadge(selectedWorkOrder.status)}</p>
                </Col>
              </Row>

              <Row className="mb-3">
                <Col md={6}>
                  <strong>Booking:</strong>
                  <p>{getBookingInfo(selectedWorkOrder.booking_id)}</p>
                </Col>
                <Col md={6}>
                  <strong>Vehicle:</strong>
                  <p>{getVehicleInfo(selectedWorkOrder.vehicle_id)}</p>
                </Col>
              </Row>

              <Row className="mb-3">
                <Col md={6}>
                  <strong>Mechanic:</strong>
                  <p>{getMechanicName(selectedWorkOrder.mechanic_id)}</p>
                </Col>
                <Col md={6}>
                  <strong>Created:</strong>
                  <p>{new Date(selectedWorkOrder.created_at).toLocaleString()}</p>
                </Col>
              </Row>

              {selectedWorkOrder.notes && (
                <Row className="mb-3">
                  <Col>
                    <strong>Notes:</strong>
                    <p>{selectedWorkOrder.notes}</p>
                  </Col>
                </Row>
              )}

              <hr />

              <h5 className="mb-3">Services</h5>
              {selectedWorkOrder.services?.length > 0 ? (
                <Table bordered size="sm">
                  <thead>
                    <tr>
                      <th>Service</th>
                      <th>Quantity</th>
                      <th>Price</th>
                      <th>Total</th>
                    </tr>
                  </thead>
                  <tbody>
                    {selectedWorkOrder.services.map((svc) => (
                      <tr key={svc.id}>
                        <td>{svc.service_name}</td>
                        <td>{svc.quantity}</td>
                        <td>${svc.price.toFixed(2)}</td>
                        <td>${(svc.price * svc.quantity).toFixed(2)}</td>
                      </tr>
                    ))}
                  </tbody>
                </Table>
              ) : (
                <p className="text-muted">No services added</p>
              )}

              <h5 className="mb-3 mt-4">Parts</h5>
              {selectedWorkOrder.parts?.length > 0 ? (
                <Table bordered size="sm">
                  <thead>
                    <tr>
                      <th>Part</th>
                      <th>Quantity</th>
                      <th>Price</th>
                      <th>Total</th>
                    </tr>
                  </thead>
                  <tbody>
                    {selectedWorkOrder.parts.map((part) => (
                      <tr key={part.id}>
                        <td>{part.sparepart_id}</td>
                        <td>{part.quantity}</td>
                        <td>${part.price.toFixed(2)}</td>
                        <td>${(part.price * part.quantity).toFixed(2)}</td>
                      </tr>
                    ))}
                  </tbody>
                </Table>
              ) : (
                <p className="text-muted">No parts added</p>
              )}

              <hr />

              <Row>
                <Col className="text-end">
                  <h4>Total: ${calculateTotal(selectedWorkOrder).toFixed(2)}</h4>
                </Col>
              </Row>

              <div className="mt-4">
                <strong>Update Status:</strong>
                <div className="btn-group w-100 mt-2" role="group">
                  <Button
                    variant={selectedWorkOrder.status === 'pending' ? 'warning' : 'outline-warning'}
                    onClick={() => handleUpdateStatus(selectedWorkOrder.id, 'pending')}
                    disabled={selectedWorkOrder.status === 'pending'}
                  >
                    Pending
                  </Button>
                  <Button
                    variant={selectedWorkOrder.status === 'in_progress' ? 'info' : 'outline-info'}
                    onClick={() => handleUpdateStatus(selectedWorkOrder.id, 'in_progress')}
                    disabled={selectedWorkOrder.status === 'in_progress'}
                  >
                    In Progress
                  </Button>
                  <Button
                    variant={selectedWorkOrder.status === 'completed' ? 'success' : 'outline-success'}
                    onClick={() => handleUpdateStatus(selectedWorkOrder.id, 'completed')}
                    disabled={selectedWorkOrder.status === 'completed'}
                  >
                    Completed
                  </Button>
                </div>
              </div>
            </>
          )}
        </Modal.Body>
        <Modal.Footer>
          <Button variant="secondary" onClick={() => setShowDetailModal(false)}>
            Close
          </Button>
        </Modal.Footer>
      </Modal>

      <Modal show={showAssignModal} onHide={() => setShowAssignModal(false)}>
        <Modal.Header closeButton>
          <Modal.Title>Assign Mechanic</Modal.Title>
        </Modal.Header>
        <Form onSubmit={handleAssignMechanic}>
          <Modal.Body>
            <Form.Group>
              <Form.Label>Select Mechanic</Form.Label>
              <Form.Select
                value={assignMechanicData.mechanic_id}
                onChange={(e) => setAssignMechanicData({ mechanic_id: e.target.value })}
                required
              >
                <option value="">Choose a mechanic</option>
                {mechanics.map(mechanic => (
                  <option key={mechanic.id} value={mechanic.id}>
                    {mechanic.name} - {mechanic.email}
                  </option>
                ))}
              </Form.Select>
            </Form.Group>
          </Modal.Body>
          <Modal.Footer>
            <Button variant="secondary" onClick={() => setShowAssignModal(false)}>
              Cancel
            </Button>
            <Button variant="primary" type="submit">
              Assign
            </Button>
          </Modal.Footer>
        </Form>
      </Modal>
    </div>
  )
}

export default WorkOrders
