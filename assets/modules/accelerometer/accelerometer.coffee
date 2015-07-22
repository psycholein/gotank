window.Accelerometer = class Accelerometer extends window.BasicModule
  config: =>
    @position = "right"

  router: (event) =>
    switch event.Task
      when "accelerometer"
        $module = $(@selector).filter("[data-name='#{event.Name}']")
        for value of event.Data
          $module.find(".#{value}").text(event.Data[value])
