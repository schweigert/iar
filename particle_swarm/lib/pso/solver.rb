require 'thread'
require 'pso/zero_vector'
require 'pso/functions/rastrigin'
require 'pso/functions/schwefel'
require 'pry'
require 'gnuplot'

module Pso
  class Solver
    def initialize(din: 5, density: 50, f: Pso::Rastrigin, center: ZeroVector[0,0,0,0,0], radius: 5.12, method: :min_by)
      @f = f.new
      @din = din
      @center = center
      @radius = radius
      @method = method
      @density = density

      generate_swarm
    end

    def generate_swarm
      Array.new(@density)
      @swarm = Array.new(@density) { generate_random_particle }
      @swarm_best = @swarm.map { |particle| [@f.f(particle), particle] }
      @swarm_speed = @swarm.map { generate_random_particle }
    end

    def generate_random_noise_particle
      @center.map { rand * 2 - 1 }
    end

    def generate_random_particle
      @center + (generate_random_noise_particle * (@radius * rand))
    end

    def perfect_particle
      if @method == :min_by
        @swarm.min_by do |element|
          @f.f(element)
        end
      else
        @swarm.max_by do |element|
          @f.f(element)
        end
      end
    end

    def solve(precision: 200000, threads: 1)
      Gnuplot.open do |gp|
        Gnuplot::Plot.new( gp ) do |plot|
          plot.terminal "gif"
          plot.output "#{self.class}.gif"
          plot.xrange "[0:200000]"
          plot.title  "Convergence"
          plot.ylabel "Minimum point"
          plot.xlabel "Interation"

          x = (0..precision).collect { |v| v.to_f }
          y = []
          Array.new(threads).map do
            Thread.new do
              ((precision / @swarm.size) / threads).times do |interation_number|
                for index in 0...@density
                  perfect = perfect_particle
                  y << @f.f(perfect)
                  new_vector = normalize(interate(@swarm[index], @swarm_best[index].last, perfect, @swarm_speed[index]))
                  @swarm_best[index] = [@f.f(new_vector), new_vector] if is_best(@swarm_best[index].first, @f.f(new_vector))
                  @swarm_speed[index] = (new_vector - @swarm[index]).normalize
                  @swarm[index] = new_vector
                end
              end
            end
          end.each do |thread|
            thread.join
          end

      plot.data = [
        Gnuplot::DataSet.new( [x, y] ) { |ds|
          ds.with = "linespoints"
          ds.title = "Array data"
        }
      ]

    end
  end

      perfect = perfect_particle
      [@f.f(perfect), perfect]
    end

    private

    def is_best(best, now)
      if @method == :min_by
        now < best
      else
        now > best
      end
    end

    def normalize(vector)
      if (vector - @center).magnitude > @radius
        return ((vector - @center).normalize * @radius) + @center
      end

      vector
    end

    def interate(vector, best, perfect, speed)
      if vector == perfect
        new_vec = vector + (best - vector).normalize * 0.2 + (generate_random_noise_particle) * rand * 0.05 + speed * 0.05
        minimal = @f.f(vector) > @f.f(new_vec)
        return @method == :min_by ? (minimal ? new_vec : vector) : ( minimal ? vector : new_vec)
      end
      vector + (perfect - vector).normalize * rand + speed * 0
    end
  end

  class SchewefelSolver < Solver
    def initialize
      super(din: 5, density: 50, f: Pso::Schwefel, center: ZeroVector[400,400,400,400,400], radius: 200.0, method: :min_by)
    end
  end
end
