window.Webcam = class Webcam extends window.BasicModule
  afterInit: =>
    img = "<image src='http://192.168.0.226:8080/?action=stream' />"
    $module = $(@selector).filter("[data-name='#{@name}']")
    $module.find('.img').html(img)

  config: =>
    @position = "left"
